package repository

import (
	"context"
	"time"
)

type DashboardData struct {
	InactiveLoans []InactiveLoan `json:"inactive_loans"`
	RecentPayments []Payment `json:"recent_payments"`
	WidgetData []Widget `json:"widget_data"`
}

type ClientDashboardResponse struct {
	ID            uint32            `json:"id"`
	FullName          string            `json:"full_name"`
	PhoneNumber   string            `json:"phone_number"`
	IdNumber      string            `json:"id_number"`
	Dob           string            `json:"dob"`
	Gender        string            `json:"gender"`
	Active        bool              `json:"active"`
	BranchName    string            `json:"branch_name"`
	AssignedStaff UserDashboardResponse `json:"assigned_staff"` 
	Overpayment   float64           `json:"overpayment"`
	DueAmount float64 `json:"due_amount"`
	CreatedBy     UserDashboardResponse `json:"created_by"` 
	CreatedAt     time.Time         `json:"created_at"`
}

type UserDashboardResponse struct {
	ID          uint32    `json:"id"`
	Fullname    string    `json:"fullname"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	BranchName  string    `json:"branch_name"`
	CreatedAt   time.Time `json:"created_at"`
	// RefreshToken string    `json:"refresh_token"`
}

type InactiveLoan struct {
	ID uint32 `json:"id"`
	Amount float64 `json:"amount"`
	RepayAmount float64 `json:"repayAmount"`
	ClientName string `json:"clientName"`
	ApprovedByName string `json:"approvedByName"`
	Client ClientDashboardResponse `json:"client"`
	LoanOfficer UserDashboardResponse `json:"loanOfficer"`
	ApprovedBy UserDashboardResponse `json:"approvedBy"`
	ApprovedOn time.Time `json:"approvedOn"`
}
// ClientID uint32 `json:"client_id"`
// LoanOfficerID uint32 `json:"loan_officer_id"`
// ApprovedBy uint32 `json:"approvedBy"`

type Payment struct {
	ID uint32 `json:"id"`
	PayingName string `json:"payingName"`
	Amount float64 `json:"amount"`
	PaidDate time.Time `json:"paidDate"`
}

type Widget struct {
	Title string `json:"title"`
	MainAmount float64 `json:"mainAmount"`
	Active float64 `json:"active"`
	ActiveTitle string `json:"activeTitle"`
	Closed float64 `json:"closed"`
	ClosedTitle string `json:"closedTitle"`
	Currency string `json:"currency"`
}

type ProductData struct {
	ID uint32 `json:"id"`
	Name string `json:"name"`
}

type ClientData struct {
	ID uint32 `json:"id"`
	Name string `json:"name"`
}

type LoanOfficerData struct {
	ID uint32 `json:"id"`
	Name string `json:"name"`
}

type BranchData struct {
	ID uint32 	`json:"id"`
	Name string	`json:"name"`
}

type HelperRepository interface {
	// helper for dashboard
	GetDashboardData(ctx context.Context) (DashboardData, error)

	// helpers for loan form
	GetProductData(ctx context.Context) ([]ProductData, error)
	GetClientData(ctx context.Context) ([]ClientData, error)
	GetLoanOfficerData(ctx context.Context) ([]LoanOfficerData, error)
	GetBranchData(ctx context.Context) ([]BranchData,  error)
}