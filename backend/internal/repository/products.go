package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type Product struct {
	ID             uint32    `json:"id"`
	BranchID       uint32    `json:"branch_id"`
	BranchName *string 	`json:"branchName"`
	LoanAmount     float64   `json:"loan_amount"`
	RepayAmount    float64   `json:"repay_amount"`
	InterestAmount float64   `json:"interest_amount"`
	UpdatedBy      uint32    `json:"updated_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type UpdateProduct struct {
	ID             uint32   `json:"id"`
	LoanAmount     *float64 `json:"loan_amount"`
	RepayAmount    *float64 `json:"repay_amount"`
	InterestAmount *float64 `json:"interest_amount"`
	UpdatedBy      uint32   `json:"updated_by"`
}

type ProductRepository interface {
	GetAllProducts(ctx context.Context, search *string, pgData *pkg.PaginationMetadata) ([]Product, pkg.PaginationMetadata, error)
	GetProductByID(ctx context.Context, id uint32) (Product, error)
	ListProductByBranch(ctx context.Context, branchID uint32, pgData *pkg.PaginationMetadata) ([]Product, error)
	CreateProduct(ctx context.Context, product *Product) (Product, error)
	DeleteProduct(ctx context.Context, id uint32) error
}
