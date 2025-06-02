package mock

import (
	"context"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.UserRepository = (*MockUserRepository)(nil)

type MockUserRepository struct {
	mockCreateUserFunc             func(ctx context.Context, user *repository.User) (repository.User, error)
	mockGetUserByIDFunc            func(ctx context.Context, id uint32) (repository.User, error)
	mockGetUserByEmailFunc         func(ctx context.Context, email string) (repository.User, error)
	mockListUsersFunc              func(ctx context.Context, category *repository.CategorySearch, pgData *pkg.PaginationMetadata) ([]repository.User, pkg.PaginationMetadata, error)
	mockUpdateUserFunc             func(ctx context.Context, user *repository.UpdateUser) (repository.User, error)
	mockUpdateUserPasswordFunc     func(ctx context.Context, email string, password string) error
	mockCheckUserExistanceFunc     func(ctx context.Context, email string) bool
	mockGetReportUserAdminDataFunc func(ctx context.Context, filters services.ReportFilters) ([]services.UserAdminsReportData, services.UserAdminsSummary, error)
	mockGetReportUserUsersDataFunc func(ctx context.Context, id uint32, filters services.ReportFilters) (services.UserUsersReportData, error)
}

func (m *MockUserRepository) CreateUser(
	ctx context.Context,
	user *repository.User,
) (repository.User, error) {
	return m.mockCreateUserFunc(ctx, user)
}
func (m *MockUserRepository) GetUserByID(ctx context.Context, id uint32) (repository.User, error) {
	return m.mockGetUserByIDFunc(ctx, id)
}

func (m *MockUserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (repository.User, error) {
	return m.mockGetUserByEmailFunc(ctx, email)
}
func (m *MockUserRepository) ListUsers(
	ctx context.Context,
	category *repository.CategorySearch,
	pgData *pkg.PaginationMetadata,
) ([]repository.User, pkg.PaginationMetadata, error) {
	return m.mockListUsersFunc(ctx, category, pgData)
}

func (m *MockUserRepository) UpdateUser(
	ctx context.Context,
	user *repository.UpdateUser,
) (repository.User, error) {
	return m.mockUpdateUserFunc(ctx, user)
}

func (m *MockUserRepository) UpdateUserPassword(
	ctx context.Context,
	email string,
	password string,
) error {
	return m.mockUpdateUserPasswordFunc(ctx, email, password)
}
func (m *MockUserRepository) CheckUserExistance(ctx context.Context, email string) bool {
	return m.mockCheckUserExistanceFunc(ctx, email)
}

func (m *MockUserRepository) GetReportUserAdminData(
	ctx context.Context,
	filters services.ReportFilters,
) ([]services.UserAdminsReportData, services.UserAdminsSummary, error) {
	return m.mockGetReportUserAdminDataFunc(ctx, filters)
}
func (m *MockUserRepository) GetReportUserUsersData(
	ctx context.Context,
	id uint32,
	filters services.ReportFilters,
) (services.UserUsersReportData, error) {
	return m.mockGetReportUserUsersDataFunc(ctx, id, filters)
}
