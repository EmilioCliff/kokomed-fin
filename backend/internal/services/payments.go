package services

import (
	"context"
	"time"
)

type MpesaCallbackData struct {
	TransactionSource string     `json:"transaction_source"`
	TransactionID     string     `json:"transaction_id"`
	AccountNumber     string     `json:"account_number"`
	PhoneNumber       string     `json:"phone_number"`
	PayingName        string     `json:"paying_name"`
	Amount            float64    `json:"amount"`
	AssignedBy        string     `json:"assigned_by"`
	AssignedTo        *uint32    `json:"assigned_to"`
	PaidDate          *time.Time `json:"paid_date"`
}

type ManualPaymentData struct {
	NonPostedID uint32 `json:"non_posted_id"`
	ClientID    uint32 `json:"client_id"`
	AdminUserID uint32 `json:"admin_user_id"`
}

type SimulatedAction struct {
	ActionType    string  `json:"actionType"`
	Description   string  `json:"description"`
	Amount        float64 `json:"amount"`
	LoanID        *uint32 `json:"loanId"`
	InstallmentID *uint32 `json:"installmentId"`
	Severity      string  `json:"severity"`
}

type SimulationResult struct {
	PaymentID uint32            `json:"paymentId"`
	UserID    uint32            `json:"userId"`
	Actions   []SimulatedAction `json:"actions"`
}

type PaymentService interface {
	ProcessCallback(ctx context.Context, callbackData *MpesaCallbackData) (uint32, error)
	TriggerManualPayment(
		ctx context.Context,
		paymentData ManualPaymentData,
		assignedBy string,
	) (uint32, error)
	UpdatePayment(
		ctx context.Context,
		paymentID uint32,
		userID uint32,
		description string,
		paymentData *MpesaCallbackData,
	) error
	DeletePayment(ctx context.Context, paymentID uint32, clientID uint32, description string) error
	SimulateUpdatePayment(
		ctx context.Context,
		paymentID uint32,
		userID uint32,
		paymentData *MpesaCallbackData,
	) (*SimulationResult, error)
	SimulateDeletePayment(
		ctx context.Context,
		paymentID uint32,
		userID uint32,
	) (*SimulationResult, error)
}
