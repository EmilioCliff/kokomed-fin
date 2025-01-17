package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type Loan struct {
	ID                 uint32     `json:"id"`
	ProductID          uint32     `json:"product_id"`
	ClientID           uint32     `json:"client_id"`
	LoanOfficerID      uint32     `json:"loan_officer_id"`
	LoanPurpose        *string    `json:"loan_purpose"`
	DueDate            *time.Time `json:"due_date"`
	ApprovedBy         uint32     `json:"approved_by"`
	DisbursedOn        *time.Time `json:"disbursed_on"`
	DisbursedBy        *uint32    `json:"disbursed_by"`
	TotalInstallments  uint32     `json:"total_installments"`
	InstallmentsPeriod uint32     `json:"installments_period"`
	Status             string     `json:"status"`
	ProcessingFee      float64    `json:"processing_fee"`
	FeePaid            bool       `json:"fee_paid"`
	PaidAmount         float64    `json:"paid_amount"`
	UpdatedBy          *uint32    `json:"updated_by"`
	CreatedBy          uint32     `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
}

type DisburseLoan struct {
	ID          uint32    `json:"id"`
	DisbursedOn time.Time `json:"disbursed_on"`
	DisbursedBy uint32    `json:"disbursed_by"`
	DueDate     time.Time `json:"due_date"`
}

type UpdateLoan struct {
	ID         uint32  `json:"id"`
	PaidAmount float64 `json:"paid_amount"`
	UpdatedBy  *uint32 `json:"updated_by"`
}

type Installment struct {
	ID              uint32    `json:"id"`
	LoanID          uint32    `json:"loan_id"`
	InstallmentNo   uint32    `json:"installment_no"`
	Amount          float64   `json:"amount"`
	RemainingAmount float64   `json:"remaining_amount"`
	Paid            bool      `json:"paid"`
	PaidAt          time.Time `json:"paid_at"`
	DueDate         time.Time `json:"due_date"`
}

type UpdateInstallment struct {
	ID              uint32     `json:"id"`
	RemainingAmount float64    `json:"remaining_amount"`
	Paid            *bool      `json:"paid"`
	PaidAt          *time.Time `json:"paid_at"`
}

type Category struct {
	BranchID    *uint32 `json:"branch_id"`
	ClientID    *uint32 `json:"client_id"`
	LoanOfficer *uint32 `json:"loan_officer"`
	Status      *string `json:"status"`
}

type LoansRepository interface {
	CreateLoan(ctx context.Context, loan *Loan) (Loan, error)
	DisburseLoan(ctx context.Context, disburseLoan *DisburseLoan) error
	TransferLoan(ctx context.Context, officerId uint32, loanId uint32, adminId uint32) error
	GetLoanByID(ctx context.Context, id uint32) (Loan, error)
	GetClientActiceLoan(ctx context.Context, clientID uint32) (uint32, error)
	ListLoans(ctx context.Context, category *Category, pgData *pkg.PaginationMetadata) ([]Loan, error)
	DeleteLoan(ctx context.Context, id uint32) error

	// use client overpayment to pay loan

	GetLoanInstallments(ctx context.Context, id uint32, pgData *pkg.PaginationMetadata) ([]Installment, error)
	GetInstallment(ctx context.Context, id uint32) (Installment, error)
	UpdateInstallment(ctx context.Context, installment *UpdateInstallment) (Installment, error)
}
