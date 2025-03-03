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

func (p *PaymentService) ProcessCallback(ctx context.Context, callbackData *services.MpesaCallbackData) (uint32, error) {
	params := &repository.NonPosted{
		TransactionSource: callbackData.TransactionSource,
		TransactionNumber: callbackData.TransactionID,
		AccountNumber:     callbackData.AccountNumber,
		PhoneNumber:       callbackData.PhoneNumber,
		PayingName:        callbackData.PayingName,
		Amount:            callbackData.Amount,
		AssignedTo:        callbackData.AssignedTo,
		AssignedBy: 	   callbackData.AssignedBy,
		PaidDate:          time.Now(),
	}

	if callbackData.PaidDate != nil {
		params.PaidDate = *callbackData.PaidDate
	}

	loanID := uint32(0)
	// if the payment is assigned to a client create a non posted and update the loan in tx
	if params.AssignedTo != nil {
		var err error
		loanID, err = p.mySQL.Loans.GetClientActiceLoan(ctx, *params.AssignedTo)
		if err != nil {
			if pkg.ErrorCode(err) == pkg.NOT_FOUND_ERROR {
				// if client doesnt have a active loan add payment to the overpayment
				if err = p.mySQL.Clients.UpdateClientOverpayment(ctx, params.AccountNumber, params.Amount); err != nil {
					return 0, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update loan overpayment: %s", err.Error())
				}

				return 0, nil
			}

			return 0, err
		}
		
		log.Println("Active Loan: ", loanID)

		err = p.db.ExecTx(ctx, func(q *generated.Queries) error {
			loan := &repository.UpdateLoan{
				ID:         loanID,
				PaidAmount: callbackData.Amount,
			}

			// create non posted
			nonPostedParams := generated.CreateNonPostedParams{
				TransactionSource: generated.NonPostedTransactionSource(params.TransactionSource),
				TransactionNumber: params.TransactionNumber,
				AccountNumber:     params.AccountNumber,
				PhoneNumber:       params.PhoneNumber,
				PayingName:        params.PayingName,
				Amount:            params.Amount,
				PaidDate:          params.PaidDate,
				AssignedBy: 	   params.AssignedBy,
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
			return 0, err
		}
	} else {
		log.Println(params)
		_, err := p.mySQL.NonPosted.CreateNonPosted(ctx, params)
		if err != nil {
			return 0, err
		}
	}

	return loanID, nil
}

func (p *PaymentService) TriggerManualPayment(ctx context.Context, paymentData services.ManualPaymentData) (uint32, error) {
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
				Valid: true,
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

	return loanID, err
}

func helperUpdateLoan(ctx context.Context, q *generated.Queries, loan *repository.UpdateLoan) error {
	log.Println("VALIDATE START: Initial Amount: ", loan.PaidAmount)
	loanData, err := q.GetLoanPaymentData(ctx, loan.ID)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loan payment data: %s", err.Error())
	}

	// if client has any overpayment add it to help pay
	clientOverpayment, err := q.GetClientOverpayment(ctx, loanData.ClientID)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get client overpayment: %s", err.Error())
	}

	// remove the overpayment from client account
	if clientOverpayment > 0 {
		_, err := q.NullifyClientOverpayment(ctx, loanData.ClientID)
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get nullify client overpayment: %s", err.Error())
		}
	}

	originalAmount := loan.PaidAmount
	loan.PaidAmount += clientOverpayment
	log.Println("VALIDATE AFTER OVERPAYMENT: Original Amount=", originalAmount, "Overpayment=", clientOverpayment, "New Amount=", loan.PaidAmount)

	// calculate total paid do not include the processing fee
	totalPaid := 0.0

	// pay the processing fee first
	if !loanData.FeePaid {
		log.Println("Paying Processing Fee")
		if loan.PaidAmount >= loanData.ProcessingFee {
			_, err := q.UpdateLoanProcessingFeeStatus(ctx, generated.UpdateLoanProcessingFeeStatusParams{
				ID:      loan.ID,
				FeePaid: true,
			})
			if err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update loan processing fee status: %s", err.Error())
			}

			beforeFee := loan.PaidAmount
			loan.PaidAmount -= loanData.ProcessingFee
			log.Println("Processing Fee Paid: Before=", beforeFee, "Fee=", loanData.ProcessingFee, "After=", loan.PaidAmount)
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

			log.Println("Processing Fee: Wasn't enough. Added to overpayment")

			return nil
		}
	}

	installments, err := q.ListUnpaidInstallmentsByLoan(ctx, loan.ID)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list unpaid installments: %s", err.Error())
	}

	for idx, i := range installments {
		log.Println("--------------------------------------------------------")
		log.Println("In Installment: ", idx, "/", len(installments), "  of LoanID", i.LoanID)
		log.Println("BEFORE PAYMENT: loan.PaidAmount=", loan.PaidAmount, "i.RemainingAmount=", i.RemainingAmount)
		
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

			// Store current values for safety check
			oldPaidAmount := loan.PaidAmount
			remainingAmount := i.RemainingAmount
			
			// Calculate expected new amount
			expectedNewAmount := oldPaidAmount - remainingAmount
			
			// Update paid amount
			loan.PaidAmount = expectedNewAmount
			
			// Safety check
			if loan.PaidAmount != expectedNewAmount {
				log.Println("ERROR: Calculation mismatch! Expected:", expectedNewAmount, "Got:", loan.PaidAmount)
				loan.PaidAmount = expectedNewAmount
			}
			
			// Update total paid
			totalPaid += remainingAmount
			
			log.Println("AFTER PAYMENT: loan.PaidAmount=", loan.PaidAmount, "totalPaid=", totalPaid)
			log.Println("Payment Details: Old Amount=", oldPaidAmount, "Installment=", remainingAmount, "New Amount=", loan.PaidAmount)
			log.Println("Paid full installment: ", loan.PaidAmount, " Total paid: ", totalPaid, " Remaining Installment: ", remainingAmount)
		} else {
			if loan.PaidAmount == 0 {
				log.Println("Amount is zero, breaking")
				break
			}
			// pay installment partially
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
			
			// Safety check for totalPaid calculation
			if totalPaid != oldTotalPaid + partialPayment {
				log.Println("ERROR: Total Paid calculation error! Expected:", oldTotalPaid + partialPayment, "Got:", totalPaid)
				totalPaid = oldTotalPaid + partialPayment
			}
			
			oldPaidAmount := loan.PaidAmount
			loan.PaidAmount = 0
			
			log.Println("PARTIAL PAYMENT: Old Amount=", oldPaidAmount, "Partial Payment=", partialPayment, "New Remaining=", newRemainingAmount)
			log.Println("Paying Partially: ", loan.PaidAmount, " Total Paid: ", totalPaid)

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

	log.Println("Final overpayment amount:", loan.PaidAmount)
	if loan.PaidAmount > 0 {
		// if client has an overpayment for the loan link to his account
		log.Println("Adding overpayment to client account")
		_, err = q.UpdateClientOverpayment(ctx, generated.UpdateClientOverpaymentParams{
			PhoneNumber: loanData.PhoneNumber,
			Overpayment: loan.PaidAmount,
			ClientID:    loanData.ClientID,
		})
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update loan overpayment: %s", err.Error())
		}
	}

	// If no remaining installments, mark loan as completed
	installmentsLeft, err := q.ListUnpaidInstallmentsByLoan(ctx, loan.ID)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list unpaid installments: %s", err.Error())
	}

	if len(installmentsLeft) == 0 {
		log.Println("Changing status to complete")
		_, err = q.UpdateLoanStatus(ctx, generated.UpdateLoanStatusParams{
			ID:     loan.ID,
			Status: generated.LoansStatusCOMPLETED,
		})
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to change loans status: %s", err.Error())
		}
	}

	log.Println("FUNCTION COMPLETE: Final totalPaid=", totalPaid, "Final loan.PaidAmount=", loan.PaidAmount)
	return nil
}