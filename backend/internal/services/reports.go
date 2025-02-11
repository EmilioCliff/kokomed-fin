package services

import (
	"context"
	"database/sql"
	"time"
)

type ReportFilters struct {
	StartDate time.Time 
	EndDate   time.Time	
	UserId *uint32 		
	ClientId *uint32
	LoanId *uint32
}

type ClientAdminsReportData struct {
	Name           string         
	BranchName     string 
	PhoneNumber    string         
	TotalLoanGiven int64          
	DefaultedLoans int64          
	ActiveLoans    int64          
	CompletedLoans int64          
	InactiveLoans  int64          
	TotalPaid      float64
	TotalDisbursed float64
	TotalOwed      float64
	Overpayment    float64        
	RateScore      float64
	DefaultRate    float64
}

type ClientClientsReportData struct {
	Name          string                         `json:"name"`
	PhoneNumber   string                         `json:"phone_number"`
	IDNumber      sql.NullString                 `json:"id_number,omitempty"`
	Dob           sql.NullTime                   `json:"dob,omitempty"`
	BranchName    sql.NullString                 `json:"branch_name,omitempty"`
	AssignedStaff sql.NullString                 `json:"assigned_staff,omitempty"`
	Active        bool                           `json:"active"`
	Loans         []ClientClientReportDataLoans  `json:"loans"`
	Payments      []ClientClientReportDataPayments `json:"payments"`
}

type ClientClientReportDataLoans struct {
	LoanId         uint32    `json:"loan_id"`
	Status         string    `json:"status"`
	LoanAmount     float64   `json:"loan_amount"`
	RepayAmount    float64   `json:"repay_amount"`
	PaidAmount     float64   `json:"paid_amount"`
	DisbursedOn    string `json:"disbursed_on"`
	TransactionFee uint16    `json:"transaction_fee"`
	CreatedBy      string    `json:"created_by"`
	AssignedBy     string    `json:"assigned_by"`
}

type ClientClientReportDataPayments struct {
	TransactionNumber string    `json:"transaction_number"`
	TransactionSource string    `json:"transaction_source"`
	AccountNumber     string    `json:"account_number"`
	PayingName        string    `json:"paying_name"`
	AssignedBy        string    `json:"assigned_by"`
	AmountPaid        float64   `json:"amount_paid"`
	PaidDate          string `json:"paid_date"`
}

type UserAdminsReportData struct {
	FullName string `json:"full_name"`
	BranchName string `json:"branch_name"`
	Roles string `json:"roles"`
	ClientsRegistered int64 `json:"clients_registered"`
	PaymentsAssigned int64 `json:"payments_assigned"`
	ApprovedLoans int64 `json:"approved_loans"`
	ActiveLoans int64 `json:"active_loans"`
	CompletedLoans int64 `json:"completed_loans"`
	DefaultRate float64 `json:"default_rate"`
}

type UserUsersReportData struct {
	FullName string `json:"full_name"`
	Roles string `json:"roles"`
	BranchName string `json:"branch_name"`
	TotalClientsHandled int64 `json:"total_clients_handled"`
	TotalLoansApproved int64 `json:"total_loans_approved"`
	TotalLoanAmountManaged float64 `json:"total_loan_amount_managed"`
	TotalCollectedAmount float64 `json:"total_collected_amount"`
	DefaultRate float64 `json:"default_rate"`
	AssignedPayments int64 `json:"assigned_payments"`
}

type BranchReportData struct {
	BranchName string `json:"branch_name"`
	TotalClients int64 `json:"total_clients"`
	TotalUsers int64 `json:"total_users"`
	LoansIssued int64 `json:"loans_issued"`
	TotalDisbursed float64 `json:"total_disbursed"`
	TotalCollected float64 `json:"total_collected"`
	TotalOutstanding float64 `json:"total_outstanding"`
	DefaultRate float64 `json:"default_rate"`
}

type PaymentReportData struct {
	TransactionNumber string    `json:"transaction_number"`
	PayingName        string    `json:"paying_name"`
	Amount            float64   `json:"amount"`
	AccountNumber     string    `json:"account_number"`
	TransactionSource string    `json:"transaction_source"`
	PaidDate          time.Time `json:"paid_date"`
	AssignedTo	  string    `json:"assigned_to"`
	AssignedBy string 			`json:"assigned_by"`
}

type ProductReportData struct {
	ProductName string
	LoansIssued int64
	ActiveLoans int64
	CompletedLoans int64
	DefaultedLoans int64
	AmountDisbursed float64
	AmountRepaid float64
	OutstandingAmount float64
	DefaultRate float64 
}

type LoanReportData struct {
	LoanID            uint32      
	ClientName        string      
	BranchName        string      
	LoanOfficer       string      
	LoanAmount        float64     
	RepayAmount       float64     
	PaidAmount        float64     
	OutstandingAmount float64
	Status            string 
	DueDate           *time.Time
	TotalInstallments uint32      
	PaidInstallments  int64       
	DisbursedDate *time.Time
	DefaultRisk       float64
}

type LoanReportDataById struct {
	LoanID                uint32                                `json:"loan_id"`
	ClientName            string                                `json:"client_name"`
	LoanAmount            float64                               `json:"loan_amount"`
	RepayAmount           float64                               `json:"repay_amount"`
	PaidAmount            float64                               `json:"paid_amount"`
	Status                string                                `json:"status"`
	TotalInstallments     uint32                                `json:"total_installments"`
	PaidInstallments      int64                                 `json:"paid_installments"`
	RemainingInstallments uint32                                `json:"remaining_installments"`
	InstallmentDetails    []LoanReportDataByIdInstallmentDetails `json:"installment_details"`
}

type LoanReportDataByIdInstallmentDetails struct {
	InstallmentNumber int64     `json:"installment_number"`
	InstallmentAmount float64   `json:"installment_amount"`
	RemainingAmount   float64   `json:"remaining_amount"`
	DueDate           string `json:"due_date"`
	Paid              uint32      `json:"paid"`
	PaidAt            string `json:"paid_at"`
}

type ReportService interface {
	GenerateLoansReport(ctx context.Context, format string, filters ReportFilters) (error)
	GeneratePaymentsReport(ctx context.Context, format string, filters ReportFilters) (error)
	GenerateBranchesReport(ctx context.Context, format string, filters ReportFilters) (error)
	GenerateUsersReport(ctx context.Context, format string, filters ReportFilters) (error)
	GenerateClientsReport(ctx context.Context, format string, filters ReportFilters) (error)
	GenerateProductsReport(ctx context.Context, format string, filters ReportFilters) (error)
}