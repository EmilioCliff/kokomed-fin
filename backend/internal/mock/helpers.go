package mock

import (
	"context"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.HelperRepository = (*MockHelperRepository)(nil)

type MockHelperRepository struct {
	mockGetDashboardDataFunc     func(ctx context.Context) (repository.DashboardData, error)
	mockGetProductDataFunc       func(ctx context.Context) ([]repository.ProductData, error)
	mockGetClientDataFunc        func(ctx context.Context) ([]repository.ClientData, error)
	mockGetLoanOfficerDataFunc   func(ctx context.Context) ([]repository.LoanOfficerData, error)
	mockGetBranchDataFunc        func(ctx context.Context) ([]repository.BranchData, error)
	mockGetLoanDataFunc          func(ctx context.Context) ([]repository.LoanData, error)
	mockGetUserFullnameFunc      func(ctx context.Context, id uint32) (string, error)
	mockGetClientNonPaymentsFunc func(ctx context.Context, id uint32, phoneNumber string, pgData *pkg.PaginationMetadata) ([]repository.NonPostedShort, pkg.PaginationMetadata, error)
}

func (m *MockHelperRepository) GetDashboardData(
	ctx context.Context,
) (repository.DashboardData, error) {
	return m.mockGetDashboardDataFunc(ctx)
}

func (m *MockHelperRepository) GetProductData(
	ctx context.Context,
) ([]repository.ProductData, error) {
	return m.mockGetProductDataFunc(ctx)
}
func (m *MockHelperRepository) GetClientData(ctx context.Context) ([]repository.ClientData, error) {
	return m.mockGetClientDataFunc(ctx)
}

func (m *MockHelperRepository) GetLoanOfficerData(
	ctx context.Context,
) ([]repository.LoanOfficerData, error) {
	return m.mockGetLoanOfficerDataFunc(ctx)
}
func (m *MockHelperRepository) GetBranchData(ctx context.Context) ([]repository.BranchData, error) {
	return m.mockGetBranchDataFunc(ctx)
}
func (m *MockHelperRepository) GetLoanData(ctx context.Context) ([]repository.LoanData, error) {
	return m.mockGetLoanDataFunc(ctx)
}
func (m *MockHelperRepository) GetUserFullname(ctx context.Context, id uint32) (string, error) {
	return m.mockGetUserFullnameFunc(ctx, id)
}
func (m *MockHelperRepository) GetClientNonPayments(
	ctx context.Context,
	id uint32,
	phoneNumber string,
	pgData *pkg.PaginationMetadata,
) ([]repository.NonPostedShort, pkg.PaginationMetadata, error) {
	return m.mockGetClientNonPaymentsFunc(ctx, id, phoneNumber, pgData)
}
