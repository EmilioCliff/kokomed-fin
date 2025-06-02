package mock

import (
	"context"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.BranchRepository = (*MockBranchRepository)(nil)

type MockBranchRepository struct {
	mockCreateBranchFunc        func(ctx context.Context, branch *repository.Branch) (repository.Branch, error)
	mockListBranchesFunc        func(ctx context.Context, search *string, pgData *pkg.PaginationMetadata) ([]repository.Branch, pkg.PaginationMetadata, error)
	mockGetBranchByIDFunc       func(ctx context.Context, id uint32) (repository.Branch, error)
	mockUpdateBranchFunc        func(ctx context.Context, name string, id uint32) (repository.Branch, error)
	mockGetReportBranchDataFunc func(ctx context.Context, filters services.ReportFilters) ([]services.BranchReportData, services.BranchSummary, error)
}

func (m *MockBranchRepository) CreateBranch(
	ctx context.Context,
	branch *repository.Branch,
) (repository.Branch, error) {
	return m.mockCreateBranchFunc(ctx, branch)
}

func (m *MockBranchRepository) ListBranches(
	ctx context.Context,
	search *string,
	pgData *pkg.PaginationMetadata,
) ([]repository.Branch, pkg.PaginationMetadata, error) {
	return m.mockListBranchesFunc(ctx, search, pgData)
}

func (m *MockBranchRepository) GetBranchByID(
	ctx context.Context,
	id uint32,
) (repository.Branch, error) {
	return m.mockGetBranchByIDFunc(ctx, id)
}

func (m *MockBranchRepository) UpdateBranch(
	ctx context.Context,
	name string,
	id uint32,
) (repository.Branch, error) {
	return m.mockUpdateBranchFunc(ctx, name, id)
}

func (m *MockBranchRepository) GetReportBranchData(
	ctx context.Context,
	filters services.ReportFilters,
) ([]services.BranchReportData, services.BranchSummary, error) {
	return m.mockGetReportBranchDataFunc(ctx, filters)
}
