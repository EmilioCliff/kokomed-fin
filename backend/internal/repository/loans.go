package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
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

type LoanFullData struct {
	ID                 uint32     `json:"id"`
	LoanPurpose        *string    `json:"loanPurpose"`
	DueDate            *time.Time `json:"dueDate"`
	DisbursedOn        *time.Time `json:"disbursedOn"`
	TotalInstallments  uint32     `json:"totalInstallments"`
	InstallmentsPeriod uint32     `json:"installmentsPeriod"`
	Status             string     `json:"status"`
	ProcessingFee      float64    `json:"processingFee"`
	FeePaid            bool       `json:"feePaid"`
	PaidAmount         float64    `json:"paidAmount"`
	RemainingAmount float64 `json:"remainingAmount"`
	CreatedAt          time.Time  `json:"createdAt"`
	Product 		   ProductShort `json:"product"`
	Client             ClientShort	`json:"client"`
	LoanOfficer        UserShortResponse	`json:"loanOfficer"`
	ApprovedBy         UserShortResponse	`json:"approvedBy"`
	DisbursedBy        UserShortResponse	`json:"disbursedBy"`
	UpdatedBy          UserShortResponse	`json:"updatedBy"`
	CreatedBy          UserShortResponse	`json:"createdBy"`
}

type DisburseLoan struct {
	ID          uint32    `json:"id"`
	DisbursedBy uint32    `json:"disbursedBy"`
	Status *string `json:"status"`
	FeePaid *bool 	`json:"feePaid"`
	DisbursedOn *time.Time `json:"disbursedOn"`
}

type UpdateLoan struct {
	ID         uint32  `json:"id"`
	PaidAmount float64 `json:"paid_amount"`
	UpdatedBy  *uint32 `json:"updated_by"`
}

type Installment struct {
	ID              uint32    `json:"id"`
	LoanID          uint32    `json:"loanId"`
	InstallmentNo   uint32    `json:"installmentNo"`
	Amount          float64   `json:"amount"`
	RemainingAmount float64   `json:"remainingAmount"`
	Paid            bool      `json:"paid"`
	PaidAt          string `json:"paidAt"`
	DueDate         string `json:"dueDate"`
}

type UpdateInstallment struct {
	ID              uint32     `json:"id"`
	RemainingAmount float64    `json:"remaining_amount"`
	Paid            *bool      `json:"paid"`
	PaidAt          *time.Time `json:"paid_at"`
}

type Category struct {
	BranchID    *uint32 `json:"branch_id"`
	Search *string `json:"string"`
	Statuses *string 	`json:"statuses"`
}

type LoanEvent struct {
	ID       string  `json:"id"`
	LoanID uint32 `json:"loanId"`
	ClientName   string  `json:"clientName"`
	LoanAmount   float64 `json:"loanAmount"`
	Date      *string `json:"date"`   
	PaymentDue   *float64 `json:"paymentDue,omitempty"` 
	Type         string  `json:"type"` 
	AllDay bool `json:"allDay"`
	Title string `json:"title"`
}

type ExpectedPayment struct {
	LoanId uint32	`json:"loanId"`
	BranchName string	`json:"branchName"`
	ClientName string	`json:"clientName"`
	LoanOfficerName string	`json:"loanOfficerName"`
	LoanAmount float64	`json:"loanAmount"`
	RepayAmount float64	`json:"repayAmount"`
	TotalUnpaid float64	`json:"totalUnpaid"`
	DueDate string	`json:"dueDate"`
}

type LoanShort struct {
	ID	uint32 `json:"id"`
	LoanAmount	float64 `json:"loanAmount"`
	RepayAmount	float64 `json:"repayAmount"`
	DisbursedOn	string `json:"disbursedOn"`
	DueDate	string `json:"dueDate"`
	PaidAmount	float64 `json:"paidAmount"`
	Installments []Installment  `json:"installments"`
}

type UnpaidInstallmentData struct {
	InstallmentNumber	uint32	`json:"installmentNumber"`
	RemainingAmount	float64	`json:"remainingAmount"`
	DueDate	string	`json:"dueDate"`
	LoanOfficer string `json:"loanOfficer"`
	LoanId	uint32	`json:"loanId"`
	ProductName string `json:"productName"`
	ClientId	uint32	`json:"clientId"`
	FullName	string	`json:"fullName"`
	PhoneNumber	string	`json:"phoneNumber"`
	ClientBranch string 	`json:"clientBranch"`
	TotalDueAmount float64 `json:"totalDueAmount"`
}

type LoansRepository interface {
	CreateLoan(ctx context.Context, loan *Loan) (LoanFullData, error)
	DisburseLoan(ctx context.Context, disburseLoan *DisburseLoan) (uint32, error)
	TransferLoan(ctx context.Context, officerId uint32, loanId uint32, adminId uint32) error
	GetLoanByID(ctx context.Context, id uint32) (Loan, error)
	GetClientActiceLoan(ctx context.Context, clientID uint32) (uint32, error)
	ListLoans(ctx context.Context, category *Category, pgData *pkg.PaginationMetadata) ([]LoanFullData, pkg.PaginationMetadata, error)
	DeleteLoan(ctx context.Context, id uint32) error
	GetExpectedPayments(ctx context.Context, category *Category, pgData *pkg.PaginationMetadata) ([]ExpectedPayment, pkg.PaginationMetadata, error)

	// use client overpayment to pay loan

	GetLoanInstallments(ctx context.Context, id uint32) ([]Installment, error)
	GetInstallment(ctx context.Context, id uint32) (Installment, error)
	UpdateInstallment(ctx context.Context, installment *UpdateInstallment) (Installment, error)
	ListUnpaidInstallmentsData(ctx context.Context, category *Category, pgData *pkg.PaginationMetadata) ([]UnpaidInstallmentData, pkg.PaginationMetadata, error)

	GetReportLoanData(ctx context.Context, filters services.ReportFilters) ([]services.LoanReportData, services.LoanSummary, error)
	GetReportLoanByIdData(ctx context.Context,id uint32) (services.LoanReportDataById, error)

	GetLoanEvents(ctx context.Context) ([]LoanEvent, error)
}
