package repository

import "time"

type PaymentAllocation struct {
	ID            uint32     `json:"id"`
	NonPostedID   uint32     `json:"nonPostedId"`
	LoanID        *uint32    `json:"loanId"`
	InstallmentID *uint32    `json:"installmentId"`
	Amount        float64    `json:"amount"`
	DeletedAt     *time.Time `json:"deletedAt"`
	CreatedAt     time.Time  `json:"createdAt"`
}
