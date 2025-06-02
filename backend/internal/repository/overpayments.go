package repository

import "time"

type Overpayment struct {
	ID          uint32    `json:"id"`
	ClientID    uint32    `json:"client_id"`
	PaymentID   *uint32   `json:"payment_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}
