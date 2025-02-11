package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db      *Store
	queries generated.Querier
}

func NewUserRepository(db *Store) *UserRepository {
	return &UserRepository{
		db:      db,
		queries: generated.New(db.db),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *repository.User) (repository.User, error) {
	execResult, err := r.queries.CreateUser(ctx, generated.CreateUserParams{
		FullName:     user.FullName,
		PhoneNumber:  user.PhoneNumber,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: user.RefreshToken,
		Role:         generated.UsersRole(user.Role),
		BranchID:     user.BranchID,
		UpdatedBy:    user.UpdatedBy, // same
		UpdatedAt:    user.UpdatedAt,
		CreatedBy:    user.CreatedBy, // same
	})
	if err != nil {
		return repository.User{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create user: %s", err.Error())
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.User{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	user.ID = uint32(id)

	return *user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uint32) (repository.User, error) {
	user, err := r.queries.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.User{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found")
		}

		return repository.User{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get user: %s", err.Error())
	}

	return convertGeneratedUser(user), nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (repository.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.User{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found")
		}

		return repository.User{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get user: %s", err.Error())
	}

	return convertGeneratedUser(user), nil
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, email string, password string) error {
	_, err := r.queries.UpdateUserPassword(ctx, generated.UpdateUserPasswordParams{
		Email:    email,
		Password: password,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found")
		}

		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update user password: %s", err.Error())
	}

	return nil
}

func (r *UserRepository) ListUsers(ctx context.Context, category *repository.CategorySearch, pgData *pkg.PaginationMetadata) ([]repository.User, pkg.PaginationMetadata, error) {
	params := generated.ListUsersByCategoryParams{
		Column1: "",
		FullName: "",
		Email: "",
		Column4: "",
		FINDINSET: "",
		Limit:    int32(pgData.PageSize),
		Offset:   int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
	}

	params2 := generated.CountUsersByCategoryParams{
		Column1: "",
		FullName: "",
		Email: "",
		Column4: "",
		FINDINSET: "",
	}

	if category.Search != nil {
		searchValue := "%" + *category.Search + "%"
		params.Column1 = "has_search"
		params.FullName = searchValue
		params.Email = searchValue

		params2.Column1 = "has_search"
		params2.FullName = searchValue
		params2.Email = searchValue
	}

	if category.Role != nil {
		params.Column4 = "has_status"
		params2.Column4 = "has_status"		
		params.FINDINSET = *category.Role
		params2.FINDINSET = *category.Role
	}

	users, err := r.queries.ListUsersByCategory(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no users found")
		}

		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get users: %s", err.Error())
	}

	totalUsers, err := r.queries.CountUsersByCategory(ctx, params2)
	if err != nil {
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get total loans: %s", err.Error())
	}

	result := make([]repository.User, len(users))
	for idx, user := range users {
		result[idx] = repository.User{
			ID: user.ID,
			FullName: user.FullName,
			PhoneNumber: user.PhoneNumber,
			Email: user.Email,
			Role: string(user.Role),
			BranchName: &user.BranchName,
			BranchID: user.BranchID,
			CreatedAt: user.CreatedAt,
		}
	}

	return result, pkg.CreatePaginationMetadata(uint32(totalUsers), pgData.PageSize, pgData.CurrentPage), nil
}

func (r *UserRepository) CheckUserExistance(ctx context.Context, email string) bool {
	count, _ := r.queries.CheckUserExistance(ctx, email)

	return count > 0
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *repository.UpdateUser) (repository.User, error) {
	params := generated.UpdateUserParams{
		ID: user.ID,
	}

	if user.Role != nil {
		params.Role = generated.NullUsersRole{
			Valid:     true,
			UsersRole: generated.UsersRole(*user.Role),
		}
	}

	if user.BranchID != nil {
		params.BranchID = sql.NullInt32{
			Valid: true,
			Int32: int32(*user.BranchID),
		}
	}

	if user.Password != nil {
		params.Password = sql.NullString{
			Valid:  true,
			String: *user.Password,
		}
	}

	if user.RefreshToken != nil {
		params.RefreshToken = sql.NullString{
			Valid:  true,
			String: *user.RefreshToken,
		}
	}

	if user.UpdatedBy != nil {
		params.UpdatedBy = sql.NullInt32{
			Valid: true,
			Int32: int32(*user.UpdatedBy),
		}
	}

	if user.UpdatedAt != nil {
		params.UpdatedAt = sql.NullTime{
			Valid: true,
			Time:  *user.UpdatedAt,
		}
	}

	_, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return repository.User{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update user: %s", err.Error())
	}

	return r.GetUserByID(ctx, user.ID)
}

func (r *UserRepository) GetReportUserAdminData(ctx context.Context, filters services.ReportFilters) ([]services.UserAdminsReportData, error) {
	users, err := r.GetUserAdminsReportData(ctx, GetUserAdminsReportDataParams{
		StartDate: filters.StartDate,
		EndDate: filters.EndDate,
		ActiveStart: filters.StartDate,
		ActiveEnd: filters.EndDate,
		CompletedStart: filters.StartDate,
		CompletedEnd: filters.EndDate,
		DefaultedStart: filters.StartDate,
		DefaultedEnd: filters.EndDate,
		TotalStart: filters.StartDate,
		TotalEnd: filters.EndDate,
		ClientsStart: filters.StartDate,
		ClientsEnd: filters.EndDate,
		PaymentsStart: filters.StartDate,
		PaymentsEnd: filters.EndDate,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get report user admin data: %s", err.Error())
	}

	rslt := make([]services.UserAdminsReportData, len(users))

	for i, user := range users {
		rslt[i] = services.UserAdminsReportData{
			FullName: user.Name,
			Roles: string(user.Role),
			BranchName: user.BranchName.String,
			ApprovedLoans: user.ApprovedLoans,
			ActiveLoans: user.ActiveLoans,
			CompletedLoans: user.CompletedLoans,
			DefaultRate: pkg.InterfaceFloat64(user.DefaultRate),
			ClientsRegistered: user.ClientsRegistered,
			PaymentsAssigned: user.PaymentsAssigned,
		}
	}

	return rslt, nil
}
func (r *UserRepository) GetReportUserUsersData(ctx context.Context, id uint32, filters services.ReportFilters) ([]services.UserUsersReportData, error) {
	users, err := r.GetUserUsersReportData(ctx, GetUserUsersReportDataParams{
		StartDate: filters.StartDate,
		EndDate: filters.EndDate,
		ID: id,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get report user admin data: %s", err.Error())
	}

	rslt := make([]services.UserUsersReportData, len(users))

	for i, user := range users {
		rslt[i] = services.UserUsersReportData{
			FullName: user.Name,
			Roles: string(user.Role),
			BranchName: user.Branch.String,
			TotalClientsHandled: user.TotalClientsHandled,
			TotalLoansApproved: user.LoansApproved,
			TotalLoanAmountManaged: pkg.InterfaceFloat64(user.TotalLoanAmountManaged),
			TotalCollectedAmount: pkg.InterfaceFloat64(user.TotalCollectedAmount),
			DefaultRate: pkg.InterfaceFloat64(user.DefaultRate),
			AssignedPayments: user.AssignedPayments,
		}
	}

	return rslt, nil
}

func convertGeneratedUser(user generated.User) repository.User {
	return repository.User{
		ID:              uint32(user.ID),
		FullName:        user.FullName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Password:        user.Password,
		PasswordUpdated: user.PasswordUpdated,
		RefreshToken:    user.RefreshToken,
		Role:            string(user.Role),
		BranchID:        user.BranchID,
		UpdatedAt:       user.UpdatedAt,
		UpdatedBy:       uint32(user.UpdatedBy),
		CreatedBy:       uint32(user.CreatedBy),
	}
}
