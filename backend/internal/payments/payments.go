package payments

import (
	"context"
	"database/sql"
	"fmt"
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
				err = p.db.ExecTx(ctx, func(q generated.Querier) error {
					nonPostedID, err := createNonPosted(ctx, q, params)
					if err != nil {
						return err
					}

					if err = updateOverpayment(ctx, q, repository.Overpayment{
						ClientID:    *params.AssignedTo,
						Amount:      params.Amount,
						PaymentID:   &nonPostedID,
						Description: "SYSTEM LOAN PAYMENT: no active loan adding to overpayment",
					}); err != nil {
						return err
					}

					if err = createAllocation(ctx, q, repository.PaymentAllocation{
						NonPostedID:   nonPostedID,
						Amount:        params.Amount,
						LoanID:        nil,
						InstallmentID: nil,
						Description:   "SYSTEM LOAN PAYMENT: no active loan adding to overpayment",
					}); err != nil {
						return err
					}

					return nil
				})
				return 0, err
			}

			return 0, err
		}

		if err = p.db.ExecTx(ctx, func(q generated.Querier) error {
			loan := &repository.UpdateLoan{
				ID:         loanID,
				PaidAmount: callbackData.Amount,
			}

			nonPostedID, err := createNonPosted(ctx, q, params)
			if err != nil {
				return err
			}

			return processLoanPayment(ctx, q, loan, nonPostedID, *params.AssignedTo, 0, "SYSTEM LOAN PAYMENT")
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
	assignedBy string,
) (uint32, error) {
	loanID := uint32(0)
	// adminFullname, err := p.mySQL.Helpers.GetUserFullname(ctx, paymentData.AdminUserID)
	// if err != nil {
	// 	return 0, err
	// }
	err := p.db.ExecTx(ctx, func(q generated.Querier) error {
		_, err := q.AssignNonPosted(ctx, generated.AssignNonPostedParams{
			ID: paymentData.NonPostedID,
			AssignTo: sql.NullInt32{
				Valid: true,
				Int32: int32(paymentData.ClientID),
			},
			TransactionSource: generated.NonPostedTransactionSourceINTERNAL,
			AssignedBy: sql.NullString{
				Valid:  true,
				String: assignedBy,
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
					ClientID:    paymentData.ClientID,
					Amount:      nonPosted.Amount,
					PaymentID:   &paymentData.NonPostedID,
					Description: "MANUAL LOAN PAYMENT: no active loan adding to overpayment",
				}); err != nil {
					return err
				}

				if err := createAllocation(ctx, q, repository.PaymentAllocation{
					NonPostedID:   paymentData.NonPostedID,
					Amount:        nonPosted.Amount,
					LoanID:        nil,
					InstallmentID: nil,
					Description:   "MANUAL LOAN PAYMENT: no active loan adding to overpayment",
				}); err != nil {
					return err
				}

				return nil
			}

			return err
		}

		if err = processLoanPayment(ctx, q, &repository.UpdateLoan{
			ID:         loanID,
			PaidAmount: nonPosted.Amount,
			UpdatedBy:  &paymentData.AdminUserID,
		}, paymentData.NonPostedID, paymentData.ClientID, 0, "MANUAL LOAN PAYMENT"); err != nil {
			return err
		}

		return nil
	})

	return loanID, err
}

func (p *PaymentService) UpdatePayment(
	ctx context.Context,
	paymentID uint32,
	userID uint32,
	description string,
	paymentData *services.MpesaCallbackData,
) error {
	nonPosted, err := p.mySQL.NonPosted.GetNonPosted(ctx, paymentID)
	if err != nil {
		return err
	}

	if nonPosted.TransactionSource != "INTERNAL" {
		return pkg.Errorf(
			pkg.INVALID_ERROR,
			"transaction source is not 'INTERNAL'. only 'INTERNAL' transactions can be updated",
		)
	}

	descriptionUpdate := fmt.Sprintf("UPDATE LOAN PAYMENT: %s", description)
	nonPostedParams := &repository.NonPosted{
		ID:                 paymentID,
		TransactionSource:  paymentData.TransactionSource,
		TransactionNumber:  paymentData.TransactionID,
		AccountNumber:      paymentData.AccountNumber,
		PhoneNumber:        paymentData.PhoneNumber,
		PayingName:         paymentData.PayingName,
		Amount:             paymentData.Amount,
		AssignedBy:         paymentData.AssignedBy,
		PaidDate:           time.Now(),
		DeletedDescription: &descriptionUpdate,
	}
	if paymentData.AssignedTo == nil {
		nonPostedParams.AssignedTo = nil
		if err := p.mySQL.NonPosted.UpdateNonPosted(ctx, nonPostedParams); err != nil {
			return err
		}

		return nil
	}

	err = p.db.ExecTx(ctx, func(q generated.Querier) error {
		allocations, err := q.ListPaymentAllocationsByNonPostedId(ctx, paymentID)
		if err != nil {
			return err
		}

		if len(allocations) == 0 {
			return pkg.Errorf(
				pkg.INVALID_ERROR,
				"action cannot be performed. payment lacks enough data to reverse payments",
			)
		}

		revertedAmount := 0.0
		loanID := uint32(0)
		for _, allocation := range allocations {
			if allocation.InstallmentID.Valid {
				loanID = uint32(allocation.LoanID.Int32)
				revertedAmount += allocation.Amount
				if err := revertInstallment(ctx, q, uint32(allocation.InstallmentID.Int32), allocation.Amount); err != nil {
					return err
				}
			} else {
				if err := deductOverpayment(ctx, q, repository.Overpayment{
					ClientID:  *paymentData.AssignedTo,
					Amount:    allocation.Amount,
					PaymentID: &paymentID,
					CreatedBy: paymentData.AssignedBy,
					Description: fmt.Sprintf(
						"UPDATE LOAN PAYMENT: REDUCING OVERPAYMENT: %s",
						description,
					),
				}); err != nil {
					return err
				}
			}
		}

		_, err = q.DeletePaymentAllocationsByNonPostedId(
			ctx,
			generated.DeletePaymentAllocationsByNonPostedIdParams{
				NonPostedID: paymentID,
				DeletedDescription: sql.NullString{
					Valid: true,
					String: fmt.Sprintf(
						"UPDATE LOAN PAYMENT: DELETING ALLOCATION: %s",
						description,
					),
				},
			},
		)
		if err != nil {
			return pkg.Errorf(
				pkg.INTERNAL_ERROR,
				"failed to delete payment allocations: %s",
				err.Error(),
			)
		}

		if revertedAmount > 0 {
			loanStatus, err := q.GetLoanStatus(ctx, loanID)
			if err != nil && err != sql.ErrNoRows {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loan status: %s", err.Error())
			}

			if loanStatus != generated.LoansStatusACTIVE {
				hasActiveLoan, err := q.CheckActiveLoanForClient(ctx, *paymentData.AssignedTo)
				if err != nil {
					return pkg.Errorf(
						pkg.INTERNAL_ERROR,
						"failed to check if client has an active loan: %s",
						err.Error(),
					)
				}

				if hasActiveLoan {
					return pkg.Errorf(
						pkg.INVALID_ERROR,
						"loan is status will change to active and client has another active loan",
					)
				}

				_, err = q.UpdateLoanStatus(ctx, generated.UpdateLoanStatusParams{
					ID:     loanID,
					Status: generated.LoansStatusACTIVE,
				})
				if err != nil {
					return pkg.Errorf(
						pkg.INTERNAL_ERROR,
						"failed to change loans status: %s",
						err.Error(),
					)
				}
			}
		}

		nonPostedParams.AssignedTo = paymentData.AssignedTo
		if err := p.mySQL.NonPosted.UpdateNonPosted(ctx, nonPostedParams); err != nil {
			return err
		}

		loan := &repository.UpdateLoan{
			ID:         loanID,
			PaidAmount: paymentData.Amount,
		}

		return processLoanPayment(
			ctx,
			q,
			loan,
			paymentID,
			*paymentData.AssignedTo,
			revertedAmount,
			fmt.Sprintf("UPDATE LOAN PAYMENT: %s", description),
		)
	})

	return err
}

func (p *PaymentService) DeletePayment(
	ctx context.Context,
	paymentID uint32,
	userID uint32,
	description string,
) error {
	paymentData, err := p.mySQL.NonPosted.GetNonPosted(ctx, paymentID)
	if err != nil {
		return err
	}

	if paymentData.TransactionSource != "INTERNAL" {
		return pkg.Errorf(
			pkg.INVALID_ERROR,
			"transaction source is not 'INTERNAL'. only 'INTERNAL' transactions can be deleted",
		)
	}

	if paymentData.AssignedTo == nil {
		if err = p.mySQL.NonPosted.DeleteNonPosted(ctx, paymentID, fmt.Sprintf(
			"DELETE PAYMENT: DELETING PAYMENT: %s",
			description,
		)); err != nil {
			return err
		}

		return nil
	}

	err = p.db.ExecTx(ctx, func(q generated.Querier) error {
		allocations, err := q.ListPaymentAllocationsByNonPostedId(ctx, paymentID)
		if err != nil {
			return err
		}

		if len(allocations) == 0 {
			return pkg.Errorf(
				pkg.INVALID_ERROR,
				"action cannot be performed. payment lacks enough data to reverse payments",
			)
		}

		revertedAmount := 0.0
		loanID := uint32(0)
		for _, allocation := range allocations {
			if allocation.InstallmentID.Valid {
				loanID = uint32(allocation.LoanID.Int32)
				revertedAmount += allocation.Amount
				if err := revertInstallment(ctx, q, uint32(allocation.InstallmentID.Int32), allocation.Amount); err != nil {
					return err
				}
			} else {
				if err := deductOverpayment(ctx, q, repository.Overpayment{
					ClientID:  *paymentData.AssignedTo,
					Amount:    allocation.Amount,
					PaymentID: &paymentID,
					CreatedBy: paymentData.AssignedBy,
					Description: fmt.Sprintf(
						"DELETE PAYMENT: REDUCING OVERPAYMENT: %s",
						description,
					),
				}); err != nil {
					return err
				}
			}
		}

		_, err = q.DeletePaymentAllocationsByNonPostedId(
			ctx,
			generated.DeletePaymentAllocationsByNonPostedIdParams{
				NonPostedID: paymentID,
				DeletedDescription: sql.NullString{
					Valid: true,
					String: fmt.Sprintf(
						"DELETE PAYMENT: DELETING ALLOCATIONS: %s",
						description,
					),
				},
			},
		)
		if err != nil {
			return pkg.Errorf(
				pkg.INTERNAL_ERROR,
				"failed to delete payment allocations: %s",
				err.Error(),
			)
		}

		err = q.SoftDeleteNonPosted(ctx, generated.SoftDeleteNonPostedParams{
			ID: paymentID,
			DeletedDescription: sql.NullString{
				Valid: true,
				String: fmt.Sprintf(
					"DELETE PAYMENT: DELETING PAYMENT: %s",
					description,
				),
			},
		})
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to delete non posted: %s", err.Error())
		}

		if revertedAmount > 0 {
			loanStatus, err := q.GetLoanStatus(ctx, loanID)
			if err != nil && err != sql.ErrNoRows {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loan status: %s", err.Error())
			}

			if loanStatus != generated.LoansStatusACTIVE {
				hasActiveLoan, err := q.CheckActiveLoanForClient(ctx, *paymentData.AssignedTo)
				if err != nil {
					return pkg.Errorf(
						pkg.INTERNAL_ERROR,
						"failed to check if client has an active loan: %s",
						err.Error(),
					)
				}

				if hasActiveLoan {
					return pkg.Errorf(
						pkg.INVALID_ERROR,
						"loan status will change to active and client has another active loan",
					)
				}

			}
			_, err = q.ReduceLoan(ctx, generated.ReduceLoanParams{
				ID:         loanID,
				PaidAmount: revertedAmount,
				UpdatedBy: sql.NullInt32{
					Valid: true,
					Int32: int32(userID),
				},
			})
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update loan: %s", err.Error())
			}
		}

		return nil
	})

	return err
}
