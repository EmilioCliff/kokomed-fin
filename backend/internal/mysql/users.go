package mysql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"go.opentelemetry.io/otel/codes"
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
	tc, span := r.db.tracer.Start(ctx, "User Repo: CreateUser")
	defer span.End()

	execResult, err := r.queries.CreateUser(tc, generated.CreateUserParams{
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
		setSpanError(span, codes.Error, err, "failed to create user")
		return repository.User{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create user: %s", err.Error())
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.User{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	branch, err := r.queries.GetBranch(ctx, user.BranchID) 
	if err != nil {
		setSpanError(span, codes.Error, err, "failed to get created user branch")
		return repository.User{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get created user branch: %s", err.Error())
	}

	user.ID = uint32(id)
	user.BranchName = &branch.Name

	return *user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uint32) (repository.User, error) {
	tc, span := r.db.tracer.Start(ctx, "User Repo: GetUserByID")
	defer span.End()

	user, err := r.queries.GetUser(tc, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.User{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found")
		}

		setSpanError(span, codes.Error, err, "failed to get user")
		return repository.User{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get user: %s", err.Error())
	}

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
		BranchName: 	 &user.BranchName,
	}, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (repository.User, error) {
	tc, span := r.db.tracer.Start(ctx, "User Repo: GetUserByEmail")
	defer span.End()

	user, err := r.queries.GetUserByEmail(tc, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.User{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found")
		}

		setSpanError(span, codes.Error, err, "failed to get user")
		return repository.User{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get user: %s", err.Error())
	}

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
		BranchName: 	 &user.BranchName,
	}, nil
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, email string, password string) error {
	tc, span := r.db.tracer.Start(ctx, "User Repo: UpdateUserPassword")
	defer span.End()

	_, err := r.queries.UpdateUserPassword(tc, generated.UpdateUserPasswordParams{
		Email:    email,
		Password: password,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found")
		}

		setSpanError(span, codes.Error, err, "failed to update user password")
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update user password: %s", err.Error())
	}

	return nil
}

func (r *UserRepository) ListUsers(ctx context.Context, category *repository.CategorySearch, pgData *pkg.PaginationMetadata) ([]repository.User, pkg.PaginationMetadata, error) {
	tc, span := r.db.tracer.Start(ctx, "User Repo: ListUsers")
	defer span.End()

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

	users, err := r.queries.ListUsersByCategory(tc, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no users found")
		}

		setSpanError(span, codes.Error, err, "failed to get users")
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get users: %s", err.Error())
	}

	totalUsers, err := r.queries.CountUsersByCategory(tc, params2)
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
	tc, span := r.db.tracer.Start(ctx, "User Repo: UpdateUser")
	defer span.End()

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

	_, err := r.queries.UpdateUser(tc, params)
	if err != nil {
		setSpanError(span, codes.Error, err, "failed to update user")
		return repository.User{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update user: %s", err.Error())
	}

	return r.GetUserByID(ctx, user.ID)
}

func (r *UserRepository) GetReportUserAdminData(ctx context.Context, filters services.ReportFilters) ([]services.UserAdminsReportData, services.UserAdminsSummary, error) {
	tc, span := r.db.tracer.Start(ctx, "User Repo: GetReportUserAdminData")
	defer span.End()

	users, err := r.GetUserAdminsReportData(tc, GetUserAdminsReportDataParams{
		StartDate: filters.StartDate,
		EndDate:   filters.EndDate,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, services.UserAdminsSummary{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found")
		}

		setSpanError(span, codes.Error, err, "failed to get report user admin data")
		return nil, services.UserAdminsSummary{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get report user admin data: %s", err.Error())
	}

	rslt := make([]services.UserAdminsReportData, len(users))
	var totalClients int64
	var totalPayments int64
	var totalApprovedLoans int64
	var highestLoanApprovalUser string
	var maxLoansApproved int64

	for i, user := range users {
		approvedLoans := user.ApprovedLoans

		rslt[i] = services.UserAdminsReportData{
			FullName:          user.Name,
			Roles:             string(user.Role),
			BranchName:        user.BranchName.String,
			ApprovedLoans:     approvedLoans,
			ActiveLoans:       user.ActiveLoans,
			CompletedLoans:    user.CompletedLoans,
			DefaultRate:       pkg.InterfaceFloat64(user.DefaultRate),
			ClientsRegistered: user.ClientsRegistered,
			PaymentsAssigned:  user.PaymentsAssigned,
		}

		totalClients += user.ClientsRegistered
		totalPayments += user.PaymentsAssigned
		totalApprovedLoans += approvedLoans

		if approvedLoans > maxLoansApproved {
			maxLoansApproved = approvedLoans
			highestLoanApprovalUser = user.Name
		}
	}

	summary := services.UserAdminsSummary{
		TotalUsers:              int64(len(users)),
		TotalClients:            totalClients,
		TotalPayments:           totalPayments,
		HighestLoanApprovalUser: highestLoanApprovalUser,
	}

	return rslt, summary, nil
}

func (r *UserRepository) GetReportUserUsersData(ctx context.Context, id uint32, filters services.ReportFilters) (services.UserUsersReportData, error) {
	tc, span := r.db.tracer.Start(ctx, "User Repo: GetReportUserUsersData")
	defer span.End()

	user, err := r.GetUserUsersReportData(tc, GetUserUsersReportDataParams{
		StartDate: filters.StartDate,
		EndDate:   filters.EndDate,
		ID:        id,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return services.UserUsersReportData{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no user found")
		}

		setSpanError(span, codes.Error, err, "failed to get report user admin data")
		return services.UserUsersReportData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get report user admin data: %s", err.Error())
	}


	return convertUserReportData(user)
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

func convertUserReportData(row GetUserUsersReportDataRow) (services.UserUsersReportData, error) {
	var assignedLoans []services.UserUsersReportDataLoans
	if row.AssignedLoans != nil {
		loansByte, ok := row.AssignedLoans.([]byte)
		if !ok {
			return services.UserUsersReportData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to convert assigned loans to bytes")
		}

		err := json.Unmarshal(loansByte, &assignedLoans)
		if err != nil {
			return services.UserUsersReportData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "error unmarshalling assigned loans: %v", err)
		}
	}

	var assignedPayments []services.UserUsersReportDataPayments
	if row.AssignedPaymentsList != nil {
		paymentsByte, ok := row.AssignedPaymentsList.([]byte)
		if !ok {
			return services.UserUsersReportData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to convert assigned payments to bytes")
		}

		err := json.Unmarshal(paymentsByte, &assignedPayments)
		if err != nil {
			return services.UserUsersReportData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "error unmarshalling assigned payments: %v", err)
		}
	}

	return services.UserUsersReportData{
		Name:                   row.Name,
		Role:                   row.Role,
		Branch:                 row.Branch.String,
		TotalClientsHandled:    row.TotalClientsHandled,
		LoansApproved:          row.LoansApproved,
		TotalLoanAmountManaged: pkg.InterfaceFloat64(row.TotalLoanAmountManaged),
		TotalCollectedAmount:   pkg.InterfaceFloat64(row.TotalCollectedAmount),
		DefaultRate:            pkg.InterfaceFloat64(row.DefaultRate),
		AssignedPayments:       row.AssignedPayments,
		AssignedLoans:          assignedLoans,
		AssignedPaymentsList:   assignedPayments,
	}, nil
}
