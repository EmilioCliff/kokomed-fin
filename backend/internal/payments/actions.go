package payments

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

func (p *PaymentService) SimulateUpdatePayment(
	ctx context.Context,
	paymentID uint32,
	userID uint32,
	paymentData *services.MpesaCallbackData,
) (*services.SimulationResult, error) {
	var result services.SimulationResult
	result.PaymentID = paymentID
	result.UserID = userID

	nonPosted, err := p.mySQL.NonPosted.GetNonPosted(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if nonPosted.TransactionSource != "INTERNAL" {
		return nil, pkg.Errorf(
			pkg.INVALID_ERROR,
			"transaction source is not 'INTERNAL'. only 'INTERNAL' transactions can be updated",
		)
	}

	if paymentData.AssignedTo == nil {
		result.Actions = append(result.Actions, services.SimulatedAction{
			ActionType:    "update_payment",
			Description:   "Non-posted payment is not assigned to any client. No accounts, loans, or overpayments will be affected.",
			Amount:        paymentData.Amount,
			LoanID:        nil,
			InstallmentID: nil,
			Severity:      "info",
		})
		return &result, nil
	}

	allocations, err := p.mySQL.NonPosted.ListPaymentAllocationsByNonPostedId(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if len(allocations) == 0 {
		result.Actions = append(result.Actions, services.SimulatedAction{
			ActionType:  "not_enough_data",
			Description: "payment has no allocations. needs extra steps",
			Severity:    "danger",
		})
		return &result, nil
	}

	revertedAmount := 0.0
	var loanID uint32

	for _, allocation := range allocations {
		if allocation.InstallmentID != nil {
			revertedAmount += allocation.Amount
			loanID = *allocation.LoanID
			result.Actions = append(result.Actions, services.SimulatedAction{
				ActionType:    "revert_installments",
				Amount:        allocation.Amount,
				InstallmentID: allocation.InstallmentID,
				LoanID:        pkg.Uint32Ptr(loanID),
				Description:   "Would revert installment",
				Severity:      "warning",
			})
		} else {
			result.Actions = append(result.Actions, services.SimulatedAction{
				ActionType: "reduce_overpayment",
				Amount:     allocation.Amount,
				LoanID:     nil,
				Description: fmt.Sprintf(
					"Would deduct overpayment for client with id %d",
					*paymentData.AssignedTo,
				),
				Severity: "warning",
			})
		}
	}

	// Loan status check
	if revertedAmount > 0 {
		loanStatus, err := p.mySQL.Loans.GetLoanStatus(ctx, loanID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		if strings.ToLower(loanStatus) != "active" {
			hasActiveLoan, err := p.mySQL.Loans.GetClientActiceLoan(ctx, *paymentData.AssignedTo)
			if err != nil {
				return nil, err
			}

			if hasActiveLoan != 0 {
				result.Actions = append(result.Actions, services.SimulatedAction{
					ActionType:  "loan_status_blocked",
					Description: "Loan status change would be blocked due to existing active loan",
					Severity:    "warning",
				})
			} else {
				result.Actions = append(result.Actions, services.SimulatedAction{
					ActionType:  "loan_status_change",
					Description: "Loan status would change to ACTIVE",
					LoanID:      pkg.Uint32Ptr(loanID),
					Severity:    "warning",
				})
			}
		}
	}

	// Summary of update
	result.Actions = append(result.Actions, services.SimulatedAction{
		ActionType:  "update_payment",
		Description: "Would update non-posted ",
		Severity:    "info",
	})
	result.Actions = append(result.Actions, services.SimulatedAction{
		ActionType:  "process_loan_payment",
		Amount:      paymentData.Amount,
		Description: "Would apply new payment for non-posted",
		LoanID:      pkg.Uint32Ptr(loanID),
		Severity:    "success",
	})

	return &result, nil
}

func (p *PaymentService) SimulateDeletePayment(
	ctx context.Context,
	paymentID uint32,
	userID uint32,
) (*services.SimulationResult, error) {
	var result services.SimulationResult

	paymentData, err := p.mySQL.NonPosted.GetNonPosted(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if paymentData.TransactionSource != "INTERNAL" {
		return nil, pkg.Errorf(
			pkg.INVALID_ERROR,
			"transaction source is not 'INTERNAL'. only 'INTERNAL' transactions can be deleted",
		)
	}

	if paymentData.AssignedTo == nil {
		result.Actions = append(result.Actions, services.SimulatedAction{
			ActionType:    "delete_non_posted",
			Description:   "Non-posted payment is not assigned to any client. No accounts, loans, or overpayments will be affected.",
			Amount:        paymentData.Amount,
			LoanID:        nil,
			InstallmentID: nil,
			Severity:      "success",
		})
		return &result, nil
	}

	allocations, err := p.mySQL.NonPosted.ListPaymentAllocationsByNonPostedId(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if len(allocations) == 0 {
		result.Actions = append(result.Actions, services.SimulatedAction{
			ActionType:  "not_enough_data",
			Description: "payment has no allocations. needs extra steps",
			Severity:    "danger",
		})
		return &result, nil
	}

	revertedAmount := 0.0
	var loanID uint32

	for _, allocation := range allocations {
		if allocation.InstallmentID != nil {
			revertedAmount += allocation.Amount
			loanID = *allocation.LoanID
			result.Actions = append(result.Actions, services.SimulatedAction{
				ActionType:    "revert_installments",
				Amount:        allocation.Amount,
				InstallmentID: allocation.InstallmentID,
				LoanID:        pkg.Uint32Ptr(loanID),
				Description:   "Would revert installment",
				Severity:      "warning",
			})
		} else {
			result.Actions = append(result.Actions, services.SimulatedAction{
				ActionType: "reduce_overpayment",
				Amount:     allocation.Amount,
				LoanID:     nil,
				Description: fmt.Sprintf(
					"Would deduct overpayment for client with id %d",
					*paymentData.AssignedTo,
				),
				Severity: "warning",
			})
		}
	}

	// Check loan status only if we reverted installment(s)
	if revertedAmount > 0 {
		loanStatus, err := p.mySQL.Loans.GetLoanStatus(ctx, loanID)
		if err != nil {
			return nil, err
		}

		if strings.ToLower(loanStatus) != "active" {
			hasActiveLoan, err := p.mySQL.Loans.GetClientActiceLoan(ctx, *paymentData.AssignedTo)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}

			if hasActiveLoan != 0 {
				result.Actions = append(result.Actions, services.SimulatedAction{
					ActionType:  "loan_status_blocked",
					Description: "Loan status change would be blocked due to existing active loan",
					Severity:    "warning",
				})
			} else {
				result.Actions = append(result.Actions, services.SimulatedAction{
					ActionType:  "loan_status_change",
					Description: "Loan status would change to ACTIVE",
					LoanID:      pkg.Uint32Ptr(loanID),
					Severity:    "warning",
				})
			}
		}
	}

	// Summary of update
	result.Actions = append(result.Actions, services.SimulatedAction{
		ActionType:  "delete_non_posted",
		Description: "Would delete non-posted ",
		Severity:    "success",
	})

	return &result, nil
}
