package services

import (
	"context"
)

type MpesaCallbackData struct {
	TransactionSource string `json:"transaction_source"`
	TransactionID string  `json:"transaction_id"`
	AccountNumber string  `json:"account_number"`
	PhoneNumber   string  `json:"phone_number"`
	PayingName    string  `json:"paying_name"`
	Amount        float64 `json:"amount"`
	AssignedBy	  string `json:"assigned_by"`
	AssignedTo    *uint32 `json:"assigned_to"`
}

type ManualPaymentData struct {
	NonPostedID uint32 `json:"non_posted_id"`
	ClientID    uint32 `json:"client_id"`
	AdminUserID uint32 `json:"admin_user_id"`
}

type PaymentService interface {
	ProcessCallback(ctx context.Context, callbackData *MpesaCallbackData) (uint32, error)
	TriggerManualPayment(ctx context.Context, paymentData ManualPaymentData) (uint32, error)
}
