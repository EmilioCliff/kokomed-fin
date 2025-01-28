package repository

import "time"

type DashboardData struct {
	InactiveLoans []Loan `json:"inactive_loans"`
	RecentPayments []Payment `json:"recent_payments"`
	WidgetData []Widget `json:"widget_data"`
}

type InactiveLoan struct {
	ID uint32 `json:"id"`
	Amount float64 `json:"amount"`
	RepayAmount float64 `json:"repay_amount"`
	ClientID uint32 `json:"client_id"`
	LoanOfficerID uint32 `json:"loan_officer_id"`
	ApprovedBy uint32 `json:"approved_by"`
	ApprovedOn time.Time `json:"approved_on"`
}

type Payment struct {
	ID uint32 `json:"id"`
	BorrowerID uint32 `json:"borrower_id"`
	Amount float64 `json:"amount"`
	Date time.Time `json:"date"`
}

type Widget struct {
	Title string `json:"title"`
	MainAmount float64 `json:"main_amount"`
	Active float64 `json:"active"`
	ActiveTitle string `json:"active_title"`
	Closed float64 `json:"closed"`
	ClosedTitle string `json:"closed_title"`
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

type HelperRepository interface {
	// helper for dashboard
	GetDashboardData() (DashboardData, error)

	// helpers for loan form
	GetProductData() ([]ProductData, error)
	GetClientData() ([]ClientData, error)
	GetLoanOfficerData() ([]LoanOfficerData, error)
}