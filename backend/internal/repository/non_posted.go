package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type NonPosted struct {
	ID                uint32    `json:"id"`
	TransactionSource string    `json:"transaction_source"`
	TransactionNumber string    `json:"transaction_number"`
	AccountNumber     string    `json:"account_number"`
	PhoneNumber       string    `json:"phone_number"`
	PayingName        string    `json:"paying_name"`
	Amount            float64   `json:"amount"`
	PaidDate          time.Time `json:"paid_date"`
	AssignedTo        *uint32   `json:"assigned_to"`
}

type NonPostedCategory struct {
	Search *string `json:"search"`
	Sources *string	`json:"sources"`
}

type NonPostedRepository interface {
	CreateNonPosted(ctx context.Context, nonPosted *NonPosted) (NonPosted, error)
	GetNonPosted(ctx context.Context, id uint32) (NonPosted, error)
	GetReportPaymentData(ctx context.Context, filters services.ReportFilters) ([]services.PaymentReportData, error)
	ListNonPosted(ctx context.Context, category *NonPostedCategory, pgData *pkg.PaginationMetadata) ([]NonPosted, pkg.PaginationMetadata, error)
	ListNonPostedByTransactionSource(ctx context.Context, transactionSource string, pgData *pkg.PaginationMetadata) ([]NonPosted, error)
	ListUnassignedNonPosted(ctx context.Context, pgData *pkg.PaginationMetadata) ([]NonPosted, error)
	DeleteNonPosted(ctx context.Context, id uint32) error
}
