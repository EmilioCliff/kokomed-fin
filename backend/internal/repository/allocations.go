package repository

import "time"

type PaymentAllocation struct {
	ID                 uint32     `json:"id"`
	NonPostedID        uint32     `json:"nonPostedId"`
	LoanID             *uint32    `json:"loanId"`
	InstallmentID      *uint32    `json:"installmentId"`
	Amount             float64    `json:"amount"`
	Description        string     `json:"description"`
	DeletedAt          *time.Time `json:"deletedAt"`
	DeletedDescription *string    `json:"deletedDescription"`
	CreatedAt          time.Time  `json:"createdAt"`
}
