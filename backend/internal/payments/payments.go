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

func (p *PaymentService) ProcessCallback(ctx context.Context, callbackData *services.MpesaCallbackData) error {
	params := &repository.NonPosted{
		TransactionSource: "MPESA",
		TransactionNumber: callbackData.TransactionID,
		AccountNumber:     callbackData.AccountNumber,
		PhoneNumber:       callbackData.PhoneNumber,
		PayingName:        callbackData.PayingName,
		Amount:            callbackData.Amount,
		AssignedTo:        callbackData.AssignedTo,
		PaidDate:          time.Now(),
	}

	// if the payment is assigned to a client create a non posted and update the loan in tx
	if params.AssignedTo != nil {
		loanID, err := p.mySQL.Loans.GetClientActiceLoan(ctx, *params.AssignedTo)
		if err != nil {
			if pkg.ErrorCode(err) == pkg.NOT_FOUND_ERROR {
				// if client doesnt have a active loan add payment to the overpayment
				if err = p.mySQL.Clients.UpdateClientOverpayment(ctx, params.AccountNumber, params.Amount); err != nil {
					return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update loan overpayment: %s", err.Error())
				}

				return nil
			}

			return err
		}

		loan := &repository.UpdateLoan{
			ID:         loanID,
			PaidAmount: callbackData.Amount,
		}

		err = p.db.ExecTx(ctx, func(q *generated.Queries) error {
			// create non posted
			nonPostedParams := generated.CreateNonPostedParams{
				TransactionSource: generated.NonPostedTransactionSource(params.TransactionSource),
				TransactionNumber: params.TransactionNumber,
				AccountNumber:     params.AccountNumber,
				PhoneNumber:       params.PhoneNumber,
				PayingName:        params.PayingName,
				Amount:            params.Amount,
				PaidDate:          params.PaidDate,
			}

			if params.AssignedTo != nil {
				nonPostedParams.AssignTo = sql.NullInt32{
					Valid: true,
					Int32: int32(*params.AssignedTo),
				}
			}

			_, err := q.CreateNonPosted(ctx, nonPostedParams)
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create non posted: %s", err.Error())
			}

			// update loan
			err = helperUpdateLoan(ctx, q, loan)

			return err
		})
		if err != nil {
			return err
		}
	} else {
		_, err := p.mySQL.NonPosted.CreateNonPosted(ctx, params)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PaymentService) TriggerManualPayment(ctx context.Context, paymentData services.ManualPaymentData) error {
	err := p.db.ExecTx(ctx, func(q *generated.Queries) error {
		_, err := q.AssignNonPosted(ctx, generated.AssignNonPostedParams{
			ID: paymentData.NonPostedID,
			AssignTo: sql.NullInt32{
				Valid: true,
				Int32: int32(paymentData.ClientID),
			},
			TransactionSource: generated.NonPostedTransactionSourceINTERNAL,
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
		loanID, err := q.GetClientActiveLoan(ctx, generated.GetClientActiveLoanParams{
			ClientID: paymentData.ClientID,
			Status:   generated.LoansStatusACTIVE,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println(loanID)
				// if client doesnt have a active loan add payment to the overpayment
				_, err := q.UpdateClientOverpayment(ctx, generated.UpdateClientOverpaymentParams{
					ClientID:    paymentData.ClientID,
					Overpayment: nonPosted.Amount,
				})
				if err != nil {
					return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update client overpayment: %s", err.Error())
				}

				return nil
			}

			return err
		}

		err = helperUpdateLoan(ctx, q, &repository.UpdateLoan{
			ID:         loanID,
			PaidAmount: nonPosted.Amount,
			UpdatedBy:  &paymentData.AdminUserID,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func helperUpdateLoan(ctx context.Context, q *generated.Queries, loan *repository.UpdateLoan) error {
	amount := loan.PaidAmount
	loanData, err := q.GetLoanPaymentData(ctx, loan.ID)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loan payment data: %s", err.Error())
	}

	// calculate total paid do not include the processing fee
	totalPaid := 0.0

	// pay the processing fee first
	if !loanData.FeePaid {
		if loan.PaidAmount >= loanData.ProcessingFee {
			_, err := q.UpdateLoanProcessingFeeStatus(ctx, generated.UpdateLoanProcessingFeeStatusParams{
				ID:      loan.ID,
				FeePaid: true,
			})
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update loan processing fee status: %s", err.Error())
			}

			loan.PaidAmount -= loanData.ProcessingFee
		} else {
			// if the processing fee is not paid fully add credit to client(overpayment).
			_, err = q.UpdateClientOverpayment(ctx, generated.UpdateClientOverpaymentParams{
				PhoneNumber: loanData.PhoneNumber,
				Overpayment: loan.PaidAmount,
				ClientID:    loanData.ClientID,
			})
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update loan overpayment: %s", err.Error())
			}

			return nil
		}
	}

	installments, err := q.ListUnpaidInstallmentsByLoan(ctx, loan.ID)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list unpaid installments: %s", err.Error())
	}

	for _, i := range installments {
		if loan.PaidAmount >= i.RemainingAmount {
			// pay the installment fully
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

			loan.PaidAmount -= i.RemainingAmount
			totalPaid += i.RemainingAmount
		} else {
			// pay installment partially
			_, err := q.PayInstallment(ctx, generated.PayInstallmentParams{
				ID:              i.ID,
				RemainingAmount: i.RemainingAmount - loan.PaidAmount,
			})
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to pay installment: %s", err.Error())
			}

			totalPaid += loan.PaidAmount
			loan.PaidAmount = 0

			break // no more amount to allocate
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

	overpayment := amount - loanData.RepayAmount
	if overpayment >= 0 {
		_, err = q.UpdateLoanStatus(ctx, generated.UpdateLoanStatusParams{
			ID:     loan.ID,
			Status: generated.LoansStatusCOMPLETED,
		})
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to change loans status: %s", err.Error())
		}

		// if client has an overpayment for the loan link to his account
		if overpayment > 0 {
			_, err = q.UpdateClientOverpayment(ctx, generated.UpdateClientOverpaymentParams{
				PhoneNumber: loanData.PhoneNumber,
				Overpayment: overpayment,
				ClientID:    loanData.ClientID,
			})
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update loan overpayment: %s", err.Error())
			}
		}
	}

	return nil
}
