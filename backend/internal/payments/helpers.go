package payments

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

func processLoanPayment(
	ctx context.Context,
	q generated.Querier,
	loan *repository.UpdateLoan,
	paymentID uint32,
	clientID uint32,
	offsetAmount float64,
	allocationMessage string,
) error {
	totalPaid := 0.0

	installments, err := q.ListUnpaidInstallmentsByLoan(ctx, loan.ID)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list unpaid installments: %s", err.Error())
	}

	var allocations []repository.PaymentAllocation
	for _, i := range installments {
		if loan.PaidAmount >= i.RemainingAmount {
			_, err := q.PayInstallment(ctx, generated.PayInstallmentParams{
				ID:              i.ID,
				RemainingAmount: 0,
				Paid: sql.NullBool{
					Valid: true,
					Bool:  true,
				},
				PaidAt: sql.NullTime{
					Valid: true,
					Time:  time.Now(),
				},
			})
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to pay installment: %s", err.Error())
			}

			oldPaidAmount := loan.PaidAmount
			remainingAmount := i.RemainingAmount
			expectedNewAmount := oldPaidAmount - remainingAmount
			loan.PaidAmount = expectedNewAmount

			if loan.PaidAmount != expectedNewAmount {
				loan.PaidAmount = expectedNewAmount
			}

			totalPaid += remainingAmount
			allocations = append(allocations, repository.PaymentAllocation{
				NonPostedID:   paymentID,
				LoanID:        &loan.ID,
				InstallmentID: &i.ID,
				Amount:        remainingAmount,
				Description:   fmt.Sprintf("%s: installment paid fully", allocationMessage),
			})
		} else {
			if loan.PaidAmount == 0 {
				break
			}

			partialPayment := loan.PaidAmount
			newRemainingAmount := i.RemainingAmount - partialPayment

			_, err := q.PayInstallment(ctx, generated.PayInstallmentParams{
				ID:              i.ID,
				RemainingAmount: newRemainingAmount,
			})
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to pay installment: %s", err.Error())
			}

			oldTotalPaid := totalPaid
			totalPaid += partialPayment

			if totalPaid != oldTotalPaid+partialPayment {
				totalPaid = oldTotalPaid + partialPayment
			}

			loan.PaidAmount = 0

			allocations = append(allocations, repository.PaymentAllocation{
				NonPostedID:   paymentID,
				LoanID:        &loan.ID,
				InstallmentID: &i.ID,
				Amount:        partialPayment,
				Description:   fmt.Sprintf("%s: installment paid partially", allocationMessage),
			})

			break
		}
	}

	for _, allocation := range allocations {
		if err := createAllocation(ctx, q, allocation); err != nil {
			return err
		}
	}

	if len(installments) > 0 {
		params := generated.UpdateLoanParams{
			ID:         loan.ID,
			PaidAmount: totalPaid,
		}

		if offsetAmount > 0 {
			params.PaidAmount = totalPaid - offsetAmount
		}

		if loan.UpdatedBy != nil {
			params.UpdatedBy = sql.NullInt32{
				Valid: true,
				Int32: int32(*loan.UpdatedBy),
			}
		}

		_, err = q.UpdateLoan(ctx, params)
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update loan: %s", err.Error())
		}
	}

	if loan.PaidAmount > 0 {
		if err := updateOverpayment(ctx, q, repository.Overpayment{
			ClientID:    clientID,
			Amount:      loan.PaidAmount,
			PaymentID:   &paymentID,
			Description: fmt.Sprintf("%s: loan is paid fully adding to overpayment", allocationMessage),
		}); err != nil {
			return err
		}

		if err = createAllocation(ctx, q, repository.PaymentAllocation{
			NonPostedID:   paymentID,
			Amount:        loan.PaidAmount,
			LoanID:        nil,
			InstallmentID: nil,
			Description:   fmt.Sprintf("%s: loan is paid fully adding to overpayment", allocationMessage),
		}); err != nil {
			return err
		}
	}

	installmentsLeft, err := q.ListUnpaidInstallmentsByLoan(ctx, loan.ID)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list unpaid installments: %s", err.Error())
	}

	if len(installmentsLeft) == 0 {
		_, err = q.UpdateLoanStatus(ctx, generated.UpdateLoanStatusParams{
			ID:     loan.ID,
			Status: generated.LoansStatusCOMPLETED,
		})
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to change loans status: %s", err.Error())
		}
	}

	return nil
}

func revertInstallment(
	ctx context.Context,
	q generated.Querier,
	installmentID uint32,
	paidAmount float64,
) error {
	_, err := q.RevertInstallment(ctx, generated.RevertInstallmentParams{
		ID:              installmentID,
		RemainingAmount: paidAmount,
	})
	if err != nil {
		return pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to update installment: %s",
			err.Error(),
		)
	}

	return nil
}

func deductOverpayment(
	ctx context.Context,
	q generated.Querier,
	data repository.Overpayment,
) error {
	overpayment, err := q.GetClientOverpayment(ctx, data.ClientID)
	if err != nil {
		return pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get client overpayment: %s",
			err.Error(),
		)
	}

	if overpayment < data.Amount {
		return pkg.Errorf(
			pkg.INVALID_ERROR,
			"client overpayment is less than the amount")
	}

	_, err = q.DeductClientOverpayment(ctx, generated.DeductClientOverpaymentParams{
		ID:          data.ClientID,
		Overpayment: data.Amount,
	})
	if err != nil {
		return pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to update client overpayment: %s",
			err.Error(),
		)
	}

	params := generated.CreateClientOverpaymentTransactionParams{
		ClientID:    data.ClientID,
		Amount:      data.Amount,
		CreatedBy:   data.CreatedBy,
		Description: data.Description,
	}

	_, err = q.CreateClientOverpaymentTransaction(ctx, params)
	if err != nil {
		return pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to create client overpayment transaction: %s",
			err.Error(),
		)
	}

	return nil
}

func updateOverpayment(
	ctx context.Context,
	q generated.Querier,
	overpaymentParams repository.Overpayment,
) error {
	_, err := q.UpdateClientOverpayment(ctx, generated.UpdateClientOverpaymentParams{
		ClientID:    overpaymentParams.ClientID,
		Overpayment: overpaymentParams.Amount,
	})
	if err != nil {
		return pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to update client overpayment: %s",
			err.Error(),
		)
	}

	params := generated.CreateClientOverpaymentTransactionParams{
		ClientID: overpaymentParams.ClientID,
		Amount:   overpaymentParams.Amount,
		PaymentID: sql.NullInt32{
			Valid: true,
			Int32: int32(*overpaymentParams.PaymentID),
		},
		CreatedBy:   "SYSTEM",
		Description: overpaymentParams.Description,
	}

	_, err = q.CreateClientOverpaymentTransaction(ctx, params)
	if err != nil {
		return pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to create client overpayment transaction: %s",
			err.Error(),
		)
	}

	return nil
}

func createNonPosted(
	ctx context.Context,
	q generated.Querier,
	params *repository.NonPosted,
) (uint32, error) {
	nonPostedParams := generated.CreateNonPostedParams{
		TransactionSource: generated.NonPostedTransactionSource(params.TransactionSource),
		TransactionNumber: params.TransactionNumber,
		AccountNumber:     params.AccountNumber,
		PhoneNumber:       params.PhoneNumber,
		PayingName:        params.PayingName,
		Amount:            params.Amount,
		PaidDate:          params.PaidDate,
		AssignedBy:        params.AssignedBy,
	}

	if params.AssignedTo != nil {
		nonPostedParams.AssignTo = sql.NullInt32{
			Valid: true,
			Int32: int32(*params.AssignedTo),
		}
	}

	nonPostedExecResult, err := q.CreateNonPosted(ctx, nonPostedParams)
	if err != nil {
		return 0, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to create non posted: %s",
			err.Error(),
		)
	}

	nonPostedID, err := nonPostedExecResult.LastInsertId()
	if err != nil {
		return 0, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get last insert id: %s",
			err.Error(),
		)
	}

	return uint32(nonPostedID), nil
}

func createAllocation(
	ctx context.Context,
	q generated.Querier,
	allocationParams repository.PaymentAllocation,
) error {
	allocation := generated.CreatePaymentAllocationParams{
		NonPostedID: allocationParams.NonPostedID,
		Amount:      allocationParams.Amount,
		Description: allocationParams.Description,
	}

	if allocationParams.LoanID != nil {
		allocation.LoanID = sql.NullInt32{
			Valid: true,
			Int32: int32(*allocationParams.LoanID),
		}
	}

	if allocationParams.InstallmentID != nil {
		allocation.InstallmentID = sql.NullInt32{
			Valid: true,
			Int32: int32(*allocationParams.InstallmentID),
		}
	}

	_, err := q.CreatePaymentAllocation(ctx, allocation)
	if err != nil {
		return pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to create payment allocation: %s",
			err.Error(),
		)
	}

	return nil
}
