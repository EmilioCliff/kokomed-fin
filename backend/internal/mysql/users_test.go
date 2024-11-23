package mysql

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/mockdb"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func NewTestUserRepository() *UserRepository {
	store := NewStore(pkg.Config{})

	return NewUserRepository(store)
}

func TestUserRepository_CreateUser(t *testing.T) {
	r := NewTestUserRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.User
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(generated.CreateUserParams{FullName: "test", Password: "test"})).
					Times(1).
					Return(&mockSQLResult{lastInsertID: 1, rowsAffected: 1}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: repository.User{
				ID:       1,
				FullName: "test",
				Password: "test",
			},
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("error"))
			},
			wantErr: true,
			err:     errors.New(pkg.INTERNAL_ERROR),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.CreateUser(context.Background(), &repository.User{FullName: "test", Password: "test"})

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, errors.New(pkg.ErrorCode(err)), tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantResult, result)
			}
		})
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	r := NewTestUserRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.User
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.User{ID: 1, FullName: "test", Password: "test"}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: repository.User{
				ID:       1,
				FullName: "test",
				Password: "test",
			},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.User{}, sql.ErrNoRows)
			},
			wantErr: true,
			err:     errors.New(pkg.NOT_FOUND_ERROR),
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.User{}, errors.New("error"))
			},
			wantErr:    true,
			err:        errors.New(pkg.INTERNAL_ERROR),
			wantResult: repository.User{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.GetUserByID(context.Background(), 1)

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, errors.New(pkg.ErrorCode(err)), tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantResult, result)
			}
		})
	}
}

func TestUserRepository_GetUserPassword(t *testing.T) {
	r := NewTestUserRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult string
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("email")).
					Times(1).
					Return("test", nil)
			},
			wantErr:    false,
			err:        nil,
			wantResult: "test",
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("email")).
					Times(1).
					Return("", sql.ErrNoRows)
			},
			wantErr: true,
			err:     errors.New(pkg.NOT_FOUND_ERROR),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.GetUserPassword(context.Background(), "email")

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, errors.New(pkg.ErrorCode(err)), tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantResult, result)
			}
		})
	}
}

func TestUserRepository_ListUsers(t *testing.T) {
	r := NewTestUserRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult []repository.User
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.User{
						{ID: 1, FullName: "test", Password: "test"},
					}, nil)
			},
			wantErr:    false,
			err:        nil,
			wantResult: []repository.User{{ID: 1, FullName: "test", Password: "test"}},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			wantErr: true,
			err:     errors.New(pkg.NOT_FOUND_ERROR),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.ListUsers(context.Background(), &pkg.PaginationMetadata{})

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, errors.New(pkg.ErrorCode(err)), tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantResult, result)
			}
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	r := NewTestUserRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.User
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				// Mock UpdateUser call
				mockQueries.EXPECT().
					UpdateUser(gomock.Any(), gomock.AssignableToTypeOf(generated.UpdateUserParams{})).
					DoAndReturn(func(_ context.Context, params generated.UpdateUserParams) (sql.Result, error) {
						// Ensure timestamps match by truncating precision
						params.UpdatedAt.Time = params.UpdatedAt.Time.Truncate(time.Millisecond)
						return &mockSQLResult{lastInsertID: 1, rowsAffected: 1}, nil
					}).
					Times(1)

				// Mock GetUser call
				mockQueries.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(uint32(1))).
					Return(generated.User{ID: 1, Role: "ADMIN", Password: "testpassword"}, nil).
					Times(1)
			},
			wantErr:    false,
			wantResult: repository.User{ID: 1, Role: "ADMIN", Password: "testpassword"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			user := &repository.UpdateUser{
				ID:        1,
				Role:      pkg.StringPtr("ADMIN"),
				Password:  pkg.StringPtr("testpassword"),
				BranchID:  pkg.Uint32Ptr(10),
				UpdatedBy: pkg.Uint32Ptr(2),
				UpdatedAt: pkg.TimePtr(time.Now()),
			}

			result, err := r.UpdateUser(context.Background(), user)

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, errors.New(pkg.ErrorCode(err)), tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantResult, result)
			}
		})
	}
}
