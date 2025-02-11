package repository

import (
	"context"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type Branch struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type BranchRepository interface {
	CreateBranch(ctx context.Context, branch *Branch) (Branch, error)
	ListBranches(ctx context.Context, search *string, pgData *pkg.PaginationMetadata) ([]Branch, pkg.PaginationMetadata, error)
	GetBranchByID(ctx context.Context, id uint32) (Branch, error)
	UpdateBranch(ctx context.Context, name string, id uint32) (Branch, error)

	GetReportBranchData(ctx context.Context, filters services.ReportFilters) ([]services.BranchReportData, error)
}
