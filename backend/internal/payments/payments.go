package payments

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ services.PaymentService = (*PaymentService)(nil)

type PaymentService struct {
	mySQL *mysql.MySQLRepo
	db    *mysql.Store
}

func NewPaymentService(mySQL *mysql.MySQLRepo, store *mysql.Store) *PaymentService {
	return &PaymentService{
		mySQL: mySQL,
		db:    store,
	}
}

func (p *PaymentService) ProcessCallback(
	ctx context.Context,
	callbackData *services.MpesaCallbackData,
) (uint32, error) {
	params := &repository.NonPosted{
		TransactionSource: callbackData.TransactionSource,
		TransactionNumber: callbackData.TransactionID,
		AccountNumber:     callbackData.AccountNumber,
		PhoneNumber:       callbackData.PhoneNumber,
		PayingName:        callbackData.PayingName,
		Amount:            callbackData.Amount,
		AssignedTo:        callbackData.AssignedTo,
		AssignedBy:        callbackData.AssignedBy,
		PaidDate:          time.Now(),
	}

	if callbackData.PaidDate != nil {
		params.PaidDate = *callbackData.PaidDate
	}

	loanID := uint32(0)
	if params.AssignedTo != nil {
		var err error
		loanID, err = p.mySQL.Loans.GetClientActiceLoan(ctx, *params.AssignedTo)
		if err != nil {
			if pkg.ErrorCode(err) == pkg.NOT_FOUND_ERROR {
				err = p.db.ExecTx(ctx, func(q *generated.Queries) error {
					nonPostedID, err := createNonPosted(ctx, q, params)
					if err != nil {
						return err
					}

					if err = updateOverpayment(ctx, q, repository.Overpayment{
						ClientID:  *params.AssignedTo,
						Amount:    params.Amount,
						PaymentID: &nonPostedID,
					}); err != nil {
						return err
					}

					if err = createAllocation(ctx, q, repository.PaymentAllocation{
						NonPostedID:   nonPostedID,
						Amount:        params.Amount,
						LoanID:        nil,
						InstallmentID: nil,
					}); err != nil {
						return err
					}

					return nil
				})
				return 0, err
			}

			return 0, err
		}

		log.Println("Active Loan...: ", loanID)

		if err = p.db.ExecTx(ctx, func(q *generated.Queries) error {
			loan := &repository.UpdateLoan{
				ID:         loanID,
				PaidAmount: callbackData.Amount,
			}

			nonPostedID, err := createNonPosted(ctx, q, params)
			if err != nil {
				return err
			}

			// update loan
			err = helperUpdateLoan(ctx, q, loan, nonPostedID, *params.AssignedTo)

			// create allocations

			return err
		}); err != nil {
			return 0, err
		}
	} else {
		_, err := p.mySQL.NonPosted.CreateNonPosted(ctx, params)
		if err != nil {
			return 0, err
		}
	}

	return loanID, nil
}

func (p *PaymentService) TriggerManualPayment(
	ctx context.Context,
	paymentData services.ManualPaymentData,
) (uint32, error) {
	loanID := uint32(0)
	adminFullname, err := p.mySQL.Helpers.GetUserFullname(ctx, paymentData.AdminUserID)
	if err != nil {
		return 0, err
	}
	err = p.db.ExecTx(ctx, func(q *generated.Queries) error {
		_, err := q.AssignNonPosted(ctx, generated.AssignNonPostedParams{
			ID: paymentData.NonPostedID,
			AssignTo: sql.NullInt32{
				Valid: true,
				Int32: int32(paymentData.ClientID),
			},
			TransactionSource: generated.NonPostedTransactionSourceMPESA, // to MPESA since attaching mpesa payment to client
			AssignedBy: sql.NullString{
				Valid:  true,
				String: adminFullname,
			},
		})
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to assign non posted: %s", err.Error())
		}

		nonPosted, err := q.GetNonPosted(ctx, paymentData.NonPostedID)
		if err != nil {
			if err == sql.ErrNoRows {
				return pkg.Errorf(pkg.NOT_FOUND_ERROR, "no non posted found")
			}

			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get non posted: %s", err.Error())
		}

		// get clients active loan
		loanID, err = q.GetClientActiveLoan(ctx, generated.GetClientActiveLoanParams{
			ClientID: paymentData.ClientID,
			Status:   generated.LoansStatusACTIVE,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				if err := updateOverpayment(ctx, q, repository.Overpayment{
					ClientID:  paymentData.ClientID,
					Amount:    nonPosted.Amount,
					PaymentID: &paymentData.NonPostedID,
				}); err != nil {
					return err
				}

				return nil
			}

			return err
		}

		err = helperUpdateLoan(ctx, q, &repository.UpdateLoan{
			ID:         loanID,
			PaidAmount: nonPosted.Amount,
			UpdatedBy:  &paymentData.AdminUserID,
		}, paymentData.NonPostedID, paymentData.ClientID)
		if err != nil {
			return err
		}

		return nil
	})

	return loanID, err
}

func helperUpdateLoan(
	ctx context.Context,
	q *generated.Queries,
	loan *repository.UpdateLoan,
	paymentID uint32,
	clientID uint32,
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
			})

			break
		}
	}

	// create allocations
	for _, allocation := range allocations {
		if err := createAllocation(ctx, q, allocation); err != nil {
			return err
		}
	}

	// update the loan
	params := generated.UpdateLoanParams{
		ID:         loan.ID,
		PaidAmount: totalPaid,
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

	if loan.PaidAmount > 0 {
		if err := updateOverpayment(ctx, q, repository.Overpayment{
			ClientID:  clientID,
			Amount:    loan.PaidAmount,
			PaymentID: &paymentID,
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

func updateOverpayment(
	ctx context.Context,
	q *generated.Queries,
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
	q *generated.Queries,
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
	q *generated.Queries,
	allocationParams repository.PaymentAllocation,
) error {
	allocation := generated.CreatePaymentAllocationParams{
		NonPostedID: allocationParams.NonPostedID,
		Amount:      allocationParams.Amount,
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
