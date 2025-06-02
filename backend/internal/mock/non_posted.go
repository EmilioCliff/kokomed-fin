package mock

import (
	"context"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.NonPostedRepository = (*MockNonPostedRepository)(nil)

type MockNonPostedRepository struct {
	mockCreateNonPostedFunc                  func(ctx context.Context, nonPosted *repository.NonPosted) (repository.NonPosted, error)
	mockGetNonPostedFunc                     func(ctx context.Context, id uint32) (repository.NonPosted, error)
	mockListNonPostedFunc                    func(ctx context.Context, category *repository.NonPostedCategory, pgData *pkg.PaginationMetadata) ([]repository.NonPosted, pkg.PaginationMetadata, error)
	mockListNonPostedByTransactionSourceFunc func(ctx context.Context, transactionSource string, pgData *pkg.PaginationMetadata) ([]repository.NonPosted, error)
	mockListUnassignedNonPostedFunc          func(ctx context.Context, pgData *pkg.PaginationMetadata) ([]repository.NonPosted, error)
	mockDeleteNonPostedFunc                  func(ctx context.Context, id uint32) error
	mockGetClientNonPostedFunc               func(ctx context.Context, id uint32, phoneNumber string, pgData *pkg.PaginationMetadata) (repository.ClientNonPosted, pkg.PaginationMetadata, error)
	mockGetReportPaymentDataFunc             func(ctx context.Context, filters services.ReportFilters) ([]services.PaymentReportData, services.PaymentSummary, error)
}

func (m *MockNonPostedRepository) CreateNonPosted(
	ctx context.Context,
	nonPosted *repository.NonPosted,
) (repository.NonPosted, error) {
	return m.mockCreateNonPostedFunc(ctx, nonPosted)
}

func (m *MockNonPostedRepository) GetNonPosted(
	ctx context.Context,
	id uint32,
) (repository.NonPosted, error) {
	return m.mockGetNonPostedFunc(ctx, id)
}

func (m *MockNonPostedRepository) ListNonPosted(
	ctx context.Context,
	category *repository.NonPostedCategory,
	pgData *pkg.PaginationMetadata,
) ([]repository.NonPosted, pkg.PaginationMetadata, error) {
	return m.mockListNonPostedFunc(ctx, category, pgData)
}

func (m *MockNonPostedRepository) ListNonPostedByTransactionSource(
	ctx context.Context,
	transactionSource string,
	pgData *pkg.PaginationMetadata,
) ([]repository.NonPosted, error) {
	return m.mockListNonPostedByTransactionSourceFunc(ctx, transactionSource, pgData)
}

func (m *MockNonPostedRepository) ListUnassignedNonPosted(
	ctx context.Context,
	pgData *pkg.PaginationMetadata,
) ([]repository.NonPosted, error) {
	return m.mockListUnassignedNonPostedFunc(ctx, pgData)
}
func (m *MockNonPostedRepository) DeleteNonPosted(ctx context.Context, id uint32) error {
	return m.mockDeleteNonPostedFunc(ctx, id)
}

func (m *MockNonPostedRepository) GetClientNonPosted(
	ctx context.Context,
	id uint32,
	phoneNumber string,
	pgData *pkg.PaginationMetadata,
) (repository.ClientNonPosted, pkg.PaginationMetadata, error) {
	return m.mockGetClientNonPostedFunc(ctx, id, phoneNumber, pgData)
}

func (m *MockNonPostedRepository) GetReportPaymentData(
	ctx context.Context,
	filters services.ReportFilters,
) ([]services.PaymentReportData, services.PaymentSummary, error) {
	return m.mockGetReportPaymentDataFunc(ctx, filters)
}
