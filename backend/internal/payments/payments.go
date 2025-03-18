package payments

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ services.PaymentService = (*PaymentService)(nil)

type PaymentService struct {
	mySQL *mysql.MySQLRepo
	db    *mysql.Store
	logger *slog.Logger
	tracer trace.Tracer
}

func NewPaymentService(mySQL *mysql.MySQLRepo, store *mysql.Store) *PaymentService {
	return &PaymentService{
		mySQL: mySQL,
		db:    store,
		logger: otelslog.NewLogger("kokomed-fin"),
		tracer: otel.Tracer("kokomed-fin"),
	}
}

func (p *PaymentService) ProcessCallback(ctx context.Context, callbackData *services.MpesaCallbackData) (uint32, error) {
	tc, span := p.tracer.Start(ctx, "ProcessCallback")
	defer span.End()

	spanCtx := span.SpanContext()

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
	pLog := p.logger.With(
		slog.String("trace_id", spanCtx.TraceID().String()),
		slog.String("span_id", spanCtx.SpanID().String()),
	)

	pLog.Info("meta_data", 
		slog.String("function", "ProcessCallback"),
		slog.String("transaction_source", params.TransactionSource),
		slog.String("transaction_number", params.TransactionNumber),
		slog.String("account_number", params.AccountNumber),
		slog.String("phone_number", params.PhoneNumber),
		slog.String("paying_name", params.PayingName),	
	)

	// if the payment is assigned to a client create a non posted and update the loan in tx
	if params.AssignedTo != nil {
		var err error
		loanID, err = p.mySQL.Loans.GetClientActiceLoan(tc, *params.AssignedTo)
		if err != nil {
			// have this in a transaction
			if pkg.ErrorCode(err) == pkg.NOT_FOUND_ERROR {
				pLog.Warn("No active loan found, updating overpayment")
				
				// if client doesnt have a active loan add payment to the overpayment
				if err = p.mySQL.Clients.UpdateClientOverpayment(tc, params.AccountNumber, params.Amount); err != nil {
					pLog.Error("Failed to update loan overpayment", slog.String("error", err.Error()))
					
					setSpanError(span, codes.Error, err, "failed to update loan overpayment")
					return 0, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update loan overpayment: %s", err.Error())
				}

				// also create a non-posted for him
				_, err := p.mySQL.NonPosted.CreateNonPosted(tc, params)
				if err != nil {
					pLog.Error("Failed to create non-posted record", slog.String("error", err.Error()))
					
					setSpanError(span, codes.Error, err, "failed to create non-posted record")
					return 0, err
				}

				pLog.Info("Overpayment recorded")
				
				return 0, nil
			}

			setSpanError(span, codes.Error, err, "failed to get active loan")
			pLog.Error("Failed to get active loan", slog.String("error", err.Error()))

			return 0, err
		}

		pLog.Info("Client has active loan", slog.Uint64("loanID", uint64(loanID)))
		
		err = p.db.ExecTx(tc, func(q *generated.Queries) error {
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

			_, err := q.CreateNonPosted(tc, nonPostedParams)
			if err != nil {
				setSpanError(span, codes.Error, err, "failed to create non posted")
				pLog.Error("Failed to create non-posted record", slog.String("error", err.Error()))
				
				return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create non posted: %s", err.Error())
			}

			// update loan
			err = helperUpdateLoan(tc, q, loan, pLog)
			if err != nil {
				setSpanError(span, codes.Error, err, "failed to update loan")
				pLog.Error("Loan update failed", slog.String("error", err.Error()))
			}

			return err
		})
		if err != nil {
			pLog.Warn("transaction failed", slog.String("error", err.Error()))

			return 0, err
		}

	} else {
		pLog.Info("Creating non-posted record, no client found")
		
		_, err := p.mySQL.NonPosted.CreateNonPosted(ctx, params)
		if err != nil {
			setSpanError(span, codes.Error, err, "failed to create non-posted")
			pLog.Error("Creating non-posted failed", slog.String("error", err.Error()))
			
			return 0, err
		}
	}

	return loanID, nil
}

func (p *PaymentService) TriggerManualPayment(ctx context.Context, paymentData services.ManualPaymentData) (uint32, error) {
	tc, span := p.tracer.Start(ctx, "TriggerManualPayment")
	defer span.End()

	spanCtx := span.SpanContext()

	loanID := uint32(0)
	adminFullname, err := p.mySQL.Helpers.GetUserFullname(tc, paymentData.AdminUserID)
	if err != nil {
		return 0, err
	}

	pLog := p.logger.With(
		slog.String("trace_id", spanCtx.TraceID().String()),
		slog.String("span_id", spanCtx.SpanID().String()),
	)

	pLog.Info("meta_data", 
		slog.String("function", "TriggerManualPayment"),
		slog.Uint64("non_posted_id", uint64(paymentData.NonPostedID)),
		slog.Uint64("client_id", uint64(paymentData.ClientID)),
		slog.String("assigned_by", adminFullname),	
	)

	pLog.Info("Starting manual payment process", "non_posted_id", paymentData.NonPostedID, "client_id", paymentData.ClientID, "assigned_by", adminFullname)
	err = p.db.ExecTx(tc, func(q *generated.Queries) error {
		_, err := q.AssignNonPosted(tc, generated.AssignNonPostedParams{
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
			pLog.Warn("failed to assign non posted", slog.String("error", err.Error()))

			setSpanError(span, codes.Error, err, "failed to assign non posted")
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to assign non posted: %s", err.Error())
		}

		nonPosted, err := q.GetNonPosted(tc, paymentData.NonPostedID)
		if err != nil {
			if err == sql.ErrNoRows {
				return pkg.Errorf(pkg.NOT_FOUND_ERROR, "no non posted found")
			}

			setSpanError(span, codes.Error, err, "failed to get non posted")
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get non posted: %s", err.Error())
		}

		// get clients active loan
		loanID, err = q.GetClientActiveLoan(tc, generated.GetClientActiveLoanParams{
			ClientID: paymentData.ClientID,
			Status:   generated.LoansStatusACTIVE,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				pLog.Warn("No active loan found, updating overpayment")
				// if client doesnt have a active loan add payment to the overpayment
				_, err := q.UpdateClientOverpayment(tc, generated.UpdateClientOverpaymentParams{
					ClientID:    paymentData.ClientID,
					Overpayment: nonPosted.Amount,
				})
				if err != nil {
					setSpanError(span, codes.Error, err, "failed to update client overpayment")
					return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update client overpayment: %s", err.Error())
				}

				return nil
			}

			return err
		}

		err = helperUpdateLoan(tc, q, &repository.UpdateLoan{
			ID:         loanID,
			PaidAmount: nonPosted.Amount,
			UpdatedBy:  &paymentData.AdminUserID,
		}, pLog)
		if err != nil {
			setSpanError(span, codes.Error, err, "failed to update loan")
			return err
		}

		return nil
	})

	return loanID, err
}

func helperUpdateLoan(ctx context.Context, q *generated.Queries, loan *repository.UpdateLoan, pLog *slog.Logger) error {
	pLog.Info("Starting loan update process", "loanID", loan.ID, "initialPaidAmount", loan.PaidAmount)

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
		pLog.Info("Applying client overpayment", "clientID", loanData.ClientID, "overpayment", clientOverpayment)

		_, err := q.NullifyClientOverpayment(ctx, loanData.ClientID)
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get nullify client overpayment: %s", err.Error())
		}
	}

	originalAmount := loan.PaidAmount
	loan.PaidAmount += clientOverpayment
	pLog.Info("Adjusted loan paid amount after overpayment", "originalAmount", originalAmount, "overpayment", clientOverpayment, "newAmount", loan.PaidAmount)

	totalPaid := 0.0

	if !loanData.FeePaid {
		pLog.Info("Processing loan fee payment", "loanID", loan.ID, "feeAmount", loanData.ProcessingFee)
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
			pLog.Info("Processing Fee Paid", "beforeFee", beforeFee, "feeAmount", loanData.ProcessingFee, "afterFee", loan.PaidAmount)
		} else {
			pLog.Warn("Insufficient funds for processing fee, adding to overpayment", "remainingAmount", loan.PaidAmount)
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

	for idx, i := range installments {
		pLog.Info("Processing installment payment(full)", "installmentIndex", idx, "loan_id", i.LoanID)
		pLog.Info("BEFORE PAYMENT", "initial_amount", loan.PaidAmount, "installemt_amount", i.RemainingAmount, "total_paid", totalPaid)
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
				pLog.Warn("ERROR: Calculation mismatch!", " Expected:", expectedNewAmount, "Got:", loan.PaidAmount)
				loan.PaidAmount = expectedNewAmount
			}
			
			// Update total paid
			totalPaid += remainingAmount
			
			pLog.Info("AFTER PAYMENT", "initial_amount", oldPaidAmount, "new_amount", loan.PaidAmount, "total_paid", totalPaid, "installemt_amount", remainingAmount)
		} else {
			if loan.PaidAmount == 0 {
				pLog.Info("No amount left to allocate, stopping processing")
				break
			}

			pLog.Info("Processing installment payment(partially)", "installmentIndex", idx, "loan_id", i.LoanID)
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
			
			if totalPaid != oldTotalPaid + partialPayment {
				pLog.Warn("ERROR: Calculation mismatch!", " Expected:", oldTotalPaid + partialPayment, "Got:", totalPaid)
				totalPaid = oldTotalPaid + partialPayment
			}
			
			oldPaidAmount := loan.PaidAmount
			loan.PaidAmount = 0
			
			pLog.Info("PARTIAL PAYMENT", "initial_amount", oldPaidAmount, "new_amount", loan.PaidAmount, "total_paid", totalPaid, "installment_remaining", newRemainingAmount, "partially_paid", partialPayment)

			break 
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
		pLog.Info("Overpayment detected", "clientID", loanData.ClientID, "overpayment", loan.PaidAmount)
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
		pLog.Info("No installments left, marking loan as complete")
		_, err = q.UpdateLoanStatus(ctx, generated.UpdateLoanStatusParams{
			ID:     loan.ID,
			Status: generated.LoansStatusCOMPLETED,
		})
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to change loans status: %s", err.Error())
		}
	}

	pLog.Info("Loan update process completed", "total_paid", totalPaid, "final_paid_amount", loan.PaidAmount)
	return nil
}

func setSpanError(span trace.Span, code codes.Code, err error, description string) {
	span.RecordError(err)
	span.SetStatus(code, description)
}