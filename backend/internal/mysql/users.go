package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
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

func (r *UserRepository) ListUsers(ctx context.Context, pgData *pkg.PaginationMetadata) ([]repository.User, error) {
	users, err := r.queries.ListUsers(ctx, generated.ListUsersParams{
		Limit:  pkg.GetPageSize(),
		Offset: pkg.CalculateOffset(pgData.CurrentPage),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no users found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get users: %s", err.Error())
	}

	result := make([]repository.User, len(users))

	for i, user := range users {
		result[i] = convertGeneratedUser(user)
	}

	return result, nil
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
