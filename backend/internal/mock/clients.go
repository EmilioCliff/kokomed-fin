package mock

import (
	"context"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.ClientRepository = (*MockClientRepository)(nil)

type MockClientRepository struct {
	mockCreateClientFunc               func(ctx context.Context, client *repository.Client) (repository.ClientFullData, error)
	mockUpdateClientFunc               func(ctx context.Context, client *repository.UpdateClient) error
	mockUpdateClientOverpaymentFunc    func(ctx context.Context, phoneNumber string, overpayment float64) error
	mockListClientsFunc                func(ctx context.Context, category *repository.ClientCategorySearch, pgData *pkg.PaginationMetadata) ([]repository.ClientFullData, pkg.PaginationMetadata, error)
	mockGetClientFullDataFunc          func(ctx context.Context, clientID uint32) (repository.ClientFullData, error)
	mockGetClientIDByPhoneNumberFunc   func(ctx context.Context, phoneNumber string) (uint32, error)
	mockListClientsByBranchFunc        func(ctx context.Context, branchID uint32, pgData *pkg.PaginationMetadata) ([]repository.Client, error)
	mockListClientsByActiveStatusFunc  func(ctx context.Context, active bool, pgData *pkg.PaginationMetadata) ([]repository.Client, error)
	mockGetReportClientAdminDataFunc   func(ctx context.Context, filters services.ReportFilters) ([]services.ClientAdminsReportData, services.ClientSummary, error)
	mockGetReportClientClientsDataFunc func(ctx context.Context, id uint32, filters services.ReportFilters) (services.ClientClientsReportData, error)
}

func (m *MockClientRepository) CreateClient(
	ctx context.Context,
	client *repository.Client,
) (repository.ClientFullData, error) {
	return m.mockCreateClientFunc(ctx, client)
}

func (m *MockClientRepository) UpdateClient(
	ctx context.Context,
	client *repository.UpdateClient,
) error {
	return m.mockUpdateClientFunc(ctx, client)
}

func (m *MockClientRepository) UpdateClientOverpayment(
	ctx context.Context,
	phoneNumber string,
	overpayment float64,
) error {
	return m.mockUpdateClientOverpaymentFunc(ctx, phoneNumber, overpayment)
}
func (m *MockClientRepository) ListClients(
	ctx context.Context,
	category *repository.ClientCategorySearch,
	pgData *pkg.PaginationMetadata,
) ([]repository.ClientFullData, pkg.PaginationMetadata, error)

// GetClient(ctx context.Context, clientID uint32) (ClientFullData, error)
func (m *MockClientRepository) GetClientFullData(
	ctx context.Context,
	clientID uint32,
) (repository.ClientFullData, error)

func (m *MockClientRepository) GetClientIDByPhoneNumber(
	ctx context.Context,
	phoneNumber string,
) (uint32, error)
func (m *MockClientRepository) ListClientsByBranch(
	ctx context.Context,
	branchID uint32,
	pgData *pkg.PaginationMetadata,
) ([]repository.Client, error)
func (m *MockClientRepository) ListClientsByActiveStatus(
	ctx context.Context,
	active bool,
	pgData *pkg.PaginationMetadata,
) ([]repository.Client, error)

func (m *MockClientRepository) GetReportClientAdminData(
	ctx context.Context,
	filters services.ReportFilters,
) ([]services.ClientAdminsReportData, services.ClientSummary, error)
func (m *MockClientRepository) GetReportClientClientsData(
	ctx context.Context,
	id uint32,
	filters services.ReportFilters,
) (services.ClientClientsReportData, error)
