package services

import (
	"context"
	"time"
)

type ReportService interface {
	GenerateLoansReport(ctx context.Context, format string, filters ReportFilters) ([]byte,error)
	GeneratePaymentsReport(ctx context.Context, format string, filters ReportFilters) ([]byte, error)
	GenerateBranchesReport(ctx context.Context, format string, filters ReportFilters) ([]byte, error)
	GenerateUsersReport(ctx context.Context, format string, filters ReportFilters) ([]byte,error)
	GenerateClientsReport(ctx context.Context, format string, filters ReportFilters) ([]byte,error)
	GenerateProductsReport(ctx context.Context, format string, filters ReportFilters) ([]byte, error)
}

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

type ClientSummary struct {
	TotalClients        int64  
	NewClients          int64  
	MostActiveClient    string 
	MostLoansClient     string 
	HighestPayingClient string 
	TotalDisbursed      float64
	TotalPaid           float64
	TotalOwed           float64
}

type ClientClientsReportData struct {
	Name          string                         `json:"name"`
	PhoneNumber   string                         `json:"phone_number"`
	IDNumber      *string                 `json:"id_number,omitempty"`
	Dob           *time.Time                   `json:"dob,omitempty"`
	BranchName    string                 `json:"branch_name,omitempty"`
	AssignedStaff string                 `json:"assigned_staff,omitempty"`
	Active        string                           `json:"active"`
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
	TransactionFee uint32    `json:"transaction_fee"`
	ApprovedBy     string    `json:"approved_by"`
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
	FullName string 
	BranchName string 
	Roles string 
	ClientsRegistered int64 
	PaymentsAssigned int64 
	ApprovedLoans int64 
	ActiveLoans int64 
	CompletedLoans int64 
	DefaultRate float64 
}

type UserAdminsSummary struct {
	TotalUsers         int64   
	TotalClients       int64   
	TotalPayments      int64   
	HighestLoanApprovalUser  string  
}

type UserUsersReportData struct {
	Name                   string         `json:"name"`
	Role                   string      `json:"role"`
	Branch                 string `json:"branch"`
	TotalClientsHandled    int64          `json:"total_clients_handled"`
	LoansApproved          int64          `json:"loans_approved"`
	TotalLoanAmountManaged float64    `json:"total_loan_amount_managed"`
	TotalCollectedAmount   float64    `json:"total_collected_amount"`
	DefaultRate            float64    `json:"default_rate"`
	AssignedPayments       int64          `json:"assigned_payments"`
	AssignedLoans          []UserUsersReportDataLoans    `json:"assigned_loans"`
	AssignedPaymentsList   []UserUsersReportDataPayments    `json:"assigned_payments_list"`
}

type UserUsersReportDataLoans struct{
	LoanId	uint32	`json:"loan_id"`
	ClientName string	`json:"client_name"`
	Status	string	`json:"status"`
	LoanAmount float64	`json:"loan_amount"`
	RepayAmount float64	`json:"repay_amount"`
	PaidAmount float64	`json:"paid_amount"`
	DisbursedOn string	`json:"disbursed_on"`
}

type UserUsersReportDataPayments struct{
	TransactionNumber string	`json:"transaction_number"`
	ClientName string	`json:"client_name"`
	AmountPaid float64	`json:"amount_paid"`
	PaidDate string	`json:"paid_date"`
}

type BranchReportData struct {
	BranchName string 
	TotalClients int64 
	TotalUsers int64 
	LoansIssued int64 
	TotalDisbursed float64 
	TotalCollected float64 
	TotalOutstanding float64 
	DefaultRate float64 
}

type BranchSummary struct {
    TotalBranches         int64   
    HighestPerformingBranch string
    MostClientsBranch     string 
    TotalClients          int64   
    TotalUsers            int64   
}

type PaymentReportData struct {
	TransactionNumber string    
	PayingName        string    
	Amount            float64   
	AccountNumber     string    
	TransactionSource string    
	PaidDate          time.Time 
	AssignedTo	  string    
	AssignedBy string 			
}

type PaymentSummary struct {
	TotalPayments        int64   
	TotalAmountReceived  float64 
	MostCommonSource     string  
	MostAssignedStaff    string  
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

type ProductSummary struct {
	TotalProducts int64
    MostPopularProduct string
    MaxLoans int64
	TotalActiveLoanAmount int64
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
	DueDate           string
	TotalInstallments uint32      
	PaidInstallments  int64       
	DisbursedDate string
	DefaultRisk       float64
}

type LoanSummary struct {
	TotalLoans            int64
	TotalActiveLoans      int64
	TotalCompletedLoans   int64
	TotalDefaultedLoans   int64
	TotalDisbursedAmount  float64
	TotalRepaidAmount     float64
	TotalOutstanding      float64
	MostIssuedLoanBranch  string
	MostLoansOfficer      string
}

type LoanReportDataById struct {
	LoanID                uint32                                `json:"loan_id"`
	ClientName            string                                `json:"client_name"`
	LoanAmount            float64                               `json:"loan_amount"`
	RepayAmount           float64                               `json:"repay_amount"`
	PaidAmount            float64                               `json:"paid_amount"`
	Status                string                                `json:"status"`
	TotalInstallments     int64                                `json:"total_installments"`
	PaidInstallments      int64                                 `json:"paid_installments"`
	RemainingInstallments int64                                `json:"remaining_installments"`
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