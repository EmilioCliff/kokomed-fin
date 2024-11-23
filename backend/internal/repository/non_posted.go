package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type NonPosted struct {
	ID                uint32    `json:"id"`
	TransactionNumber string    `json:"transaction_number"`
	AccountNumber     string    `json:"account_number"`
	PhoneNumber       string    `json:"phone_number"`
	PayingName        string    `json:"paying_name"`
	Amount            float64   `json:"amount"`
	PaidDate          time.Time `json:"paid_date"`
	AssignedTo        *uint32   `json:"assigned_to"`
}
type NonPostedRepository interface {
	CreateNonPosted(ctx context.Context, nonPosted *NonPosted) (NonPosted, error)
	GetNonPosted(ctx context.Context, id uint32) (NonPosted, error)
	AssignNonPosted(ctx context.Context, id uint32, assignedTo uint32) (NonPosted, error)
	ListNonPosted(ctx context.Context, pgData *pkg.PaginationMetadata) ([]NonPosted, error)
	ListUnassignedNonPosted(ctx context.Context, pgData *pkg.PaginationMetadata) ([]NonPosted, error)
	DeleteNonPosted(ctx context.Context, id uint32) error
}
