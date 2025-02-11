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
	Name          string         
	PhoneNumber   string         
	IDNumber      sql.NullString 
	Dob           sql.NullTime   
	BranchName    sql.NullString 
	AssignedStaff sql.NullString 
	Active        bool           
	Loans         []ClientClientReportDataLoans    
	Payments      []ClientClientReportDataPayments
}

type ClientClientReportDataLoans struct {
	LoanId uint32
	Status string
	LoanAmount float64
	RepayAmount float64
	PaidAmount float64
	DisbursedOn time.Time
	TransactionFee uint16
	CreatedBy string
	AssignedBy string
}

type ClientClientReportDataPayments struct {
	TransactionNumber string
	TransactionSource string
	AccountNumber string
	PayingName string
	AssignedBy string
	AmountPaid float64
	PaidDate time.Time
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
	LoanID                uint32      
	ClientName            string      
	LoanAmount            float64     
	RepayAmount           float64     
	PaidAmount            float64     
	Status                string 	  
	TotalInstallments     uint32      
	PaidInstallments      int64       
	RemainingInstallments uint32      
	InstallmentDetails    []LoanReportDataByIdInstallmentDetails 
}

type LoanReportDataByIdInstallmentDetails struct {
	InstallmentNumber int64
	InstallmentAmount float64
	RemainingAmount float64
	DueDate time.Time
	Paid  bool
	PaidAt time.Time
}

type ReportService interface {
	GenerateLoansReport(ctx context.Context, format string, filters ReportFilters) (error)
	GeneratePaymentsReport(ctx context.Context, format string, filters ReportFilters) (error)
	GenerateBranchesReport(ctx context.Context, format string, filters ReportFilters) (error)
	GenerateUsersReport(ctx context.Context, format string, filters ReportFilters) (error)
	GenerateClientsReport(ctx context.Context, format string, filters ReportFilters) (error)
	GenerateProductsReport(ctx context.Context, format string, filters ReportFilters) (error)
}