package mysql

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/mockdb"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func NewTestNonPostedRepository() *NonPostedRepository {
	store := NewStore(pkg.Config{})

	return NewNonPostedRepository(store)
}

func TestNonPostedRepository_TestCreateNonPosted(t *testing.T) {
	r := NewTestNonPostedRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.NonPosted
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					CreateNonPosted(gomock.Any(), gomock.Eq(generated.CreateNonPostedParams{TransactionNumber: "test"})).
					Times(1).
					Return(&mockSQLResult{lastInsertID: 1, rowsAffected: 1}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: repository.NonPosted{
				ID:                1,
				TransactionNumber: "test",
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					CreateNonPosted(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&mockSQLResult{}, errors.New("internal error"))
			},
			wantErr: true,
			err:     errors.New(pkg.INTERNAL_ERROR),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.CreateNonPosted(context.Background(), &repository.NonPosted{TransactionNumber: "test"})

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

func TestNonPostedRepository_GetNonPosted(t *testing.T) {
	r := NewTestNonPostedRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.NonPosted
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetNonPosted(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.NonPosted{ID: 1, TransactionNumber: "test"}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: repository.NonPosted{
				ID:                1,
				TransactionNumber: "test",
			},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetNonPosted(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.NonPosted{}, sql.ErrNoRows)
			},
			wantErr: true,
			err:     errors.New(pkg.NOT_FOUND_ERROR),
		},
		{
			name: "Internal Server Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetNonPosted(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.NonPosted{}, errors.New("internal error"))
			},
			wantErr: true,
			err:     errors.New(pkg.INTERNAL_ERROR),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.GetNonPosted(context.Background(), 1)

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

func TestNonPostedRepository_ListNonPosted(t *testing.T) {
	r := NewTestNonPostedRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult []repository.NonPosted
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListAllNonPosted(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.NonPosted{
						{ID: 1, TransactionNumber: "test"},
					}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: []repository.NonPosted{
				{ID: 1, TransactionNumber: "test"},
			},
		},
		{
			name: "NonPosted not found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListAllNonPosted(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.NonPosted{}, nil)
			},
			wantErr:    false,
			err:        nil,
			wantResult: []repository.NonPosted{},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListAllNonPosted(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("internal error"))
			},
			wantErr: true,
			err:     errors.New(pkg.INTERNAL_ERROR),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.ListNonPosted(context.Background(), &pkg.PaginationMetadata{})

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

func TestNonPostedRepository_ListUnassignedNonPosted(t *testing.T) {
	r := NewTestNonPostedRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult []repository.NonPosted
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListUnassignedNonPosted(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.NonPosted{
						{ID: 1, TransactionNumber: "test"},
					}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: []repository.NonPosted{
				{ID: 1, TransactionNumber: "test"},
			},
		},
		{
			name: "NonPosted not found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListUnassignedNonPosted(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.NonPosted{}, nil)
			},
			wantErr:    false,
			err:        nil,
			wantResult: []repository.NonPosted{},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListUnassignedNonPosted(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("internal error"))
			},
			wantErr: true,
			err:     errors.New(pkg.INTERNAL_ERROR),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.ListUnassignedNonPosted(context.Background(), &pkg.PaginationMetadata{})

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

func TestNonPostedRepository_DeleteNonPosted(t *testing.T) {
	r := NewTestNonPostedRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					DeleteNonPosted(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "Internal Server Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					DeleteNonPosted(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("internal error"))
			},
			wantErr: true,
			err:     errors.New(pkg.INTERNAL_ERROR),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			err := r.DeleteNonPosted(context.Background(), 1)

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, errors.New(pkg.ErrorCode(err)), tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
