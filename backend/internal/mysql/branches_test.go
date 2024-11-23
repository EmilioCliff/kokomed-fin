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

func NewTestBranchRepository() *BranchRepository {
	store := NewStore(pkg.Config{})

	return NewBranchRepository(store)
}

func TestBranchRepository_CreateBranch(t *testing.T) {
	r := NewTestBranchRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.Branch
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					CreateBranch(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&mockSQLResult{lastInsertID: 1, rowsAffected: 1}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: repository.Branch{
				ID:   1,
				Name: "test",
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					CreateBranch(gomock.Any(), gomock.Any()).
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

			branch, err := r.CreateBranch(context.Background(), &repository.Branch{Name: "test"})
			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, errors.New(pkg.ErrorCode(err)), tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantResult, branch)
			}
		})
	}
}

func TestBranchRepository_ListBranches(t *testing.T) {
	r := NewTestBranchRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult []repository.Branch
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListBranches(gomock.Any()).
					Times(1).
					Return([]generated.Branch{
						{ID: 1, Name: "test"},
					}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: []repository.Branch{
				{ID: 1, Name: "test"},
			},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListBranches(gomock.Any()).
					Times(1).
					Return([]generated.Branch{}, nil)
			},
			wantErr:    false,
			err:        nil,
			wantResult: []repository.Branch{},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListBranches(gomock.Any()).
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

			branches, err := r.ListBranches(context.Background())
			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, errors.New(pkg.ErrorCode(err)), tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantResult, branches)
			}
		})
	}
}

func TestBranchRepository_GetBranchByID(t *testing.T) {
	r := NewTestBranchRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.Branch
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetBranch(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.Branch{ID: 1, Name: "test"}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: repository.Branch{
				ID:   1,
				Name: "test",
			},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetBranch(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.Branch{}, sql.ErrNoRows)
			},
			wantErr:    true,
			err:        errors.New(pkg.NOT_FOUND_ERROR),
			wantResult: repository.Branch{},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetBranch(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.Branch{}, errors.New("error"))
			},
			wantErr: true,
			err:     errors.New(pkg.INTERNAL_ERROR),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			branch, err := r.GetBranchByID(context.Background(), 1)
			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, errors.New(pkg.ErrorCode(err)), tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantResult, branch)
			}
		})
	}
}

func TestBranchRepository_UpdateBranch(t *testing.T) {
	r := NewTestBranchRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.Branch
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					UpdateBranch(gomock.Any(), gomock.Eq(generated.UpdateBranchParams{
						ID:   1,
						Name: "UpdatedBranch",
					})).
					Times(1).
					Return(&mockSQLResult{lastInsertID: 1, rowsAffected: 1}, nil)
			},
			wantErr:    false,
			err:        nil,
			wantResult: repository.Branch{ID: 1, Name: "UpdatedBranch"},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					UpdateBranch(gomock.Any(), gomock.Eq(generated.UpdateBranchParams{
						ID:   1,
						Name: "UpdatedBranch",
					})).
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

			result, err := r.UpdateBranch(context.Background(), "UpdatedBranch", 1)

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
