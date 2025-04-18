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

func NewTestProductRepository() *ProductRepository {
	store := NewStore(pkg.Config{})

	return NewProductRepository(store)
}

func TestProductRepository_GetAllProducts(t *testing.T) {
	r := NewTestProductRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult []repository.Product
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.Product{
						{ID: 42, BranchID: 1, LoanAmount: 1000},
					}, nil)
			},
			wantErr:    false,
			err:        nil,
			wantResult: []repository.Product{{ID: 42, BranchID: 1, LoanAmount: 1000}},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			wantErr: true,
			err:     errors.New(pkg.NOT_FOUND_ERROR),
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
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

			result, err := r.GetAllProducts(context.Background(), &pkg.PaginationMetadata{})

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

func TestProductRepository_GetProductByID(t *testing.T) {
	r := NewTestProductRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.Product
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.Product{ID: 42, BranchID: 1, LoanAmount: 1000}, nil)
			},
			wantErr:    false,
			err:        nil,
			wantResult: repository.Product{ID: 42, BranchID: 1, LoanAmount: 1000},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.Product{}, sql.ErrNoRows)
			},
			wantErr: true,
			err:     errors.New(pkg.NOT_FOUND_ERROR),
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.Product{}, errors.New("internal error"))
			},
			wantErr: true,
			err:     errors.New(pkg.INTERNAL_ERROR),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.GetProductByID(context.Background(), 42)

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

func TestProductRepository_ListProductByBranch(t *testing.T) {
	r := NewTestProductRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult []repository.Product
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListProductsByBranch(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.Product{
						{ID: 42, BranchID: 1, LoanAmount: 1000},
					}, nil)
			},
			wantErr:    false,
			err:        nil,
			wantResult: []repository.Product{{ID: 42, BranchID: 1, LoanAmount: 1000}},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListProductsByBranch(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			wantErr: true,
			err:     errors.New(pkg.NOT_FOUND_ERROR),
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListProductsByBranch(gomock.Any(), gomock.Any()).
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

			result, err := r.ListProductByBranch(context.Background(), 1, &pkg.PaginationMetadata{})

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

func TestProductRepository_CreateProduct(t *testing.T) {
	r := NewTestProductRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.Product
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&mockSQLResult{lastInsertID: 1, rowsAffected: 1}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: repository.Product{
				ID:         1,
				BranchID:   1,
				LoanAmount: 1000,
			},
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
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

			product, err := r.CreateProduct(context.Background(), &repository.Product{
				BranchID:   1,
				LoanAmount: 1000,
			})

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, errors.New(pkg.ErrorCode(err)), tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantResult, product)
			}
		})
	}
}

func TestProductRepository_DeleteProduct(t *testing.T) {
	r := NewTestProductRepository()

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
					DeleteProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrNoRows)
			},
			wantErr: true,
			err:     errors.New(pkg.NOT_FOUND_ERROR),
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					DeleteProduct(gomock.Any(), gomock.Any()).
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

			err := r.DeleteProduct(context.Background(), 1)

			if tc.wantErr {
				require.Error(t, err)
				require.EqualError(t, errors.New(pkg.ErrorCode(err)), tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type mockSQLResult struct {
	lastInsertID int64
	rowsAffected int64
}

func (m *mockSQLResult) LastInsertId() (int64, error) {
	return m.lastInsertID, nil
}

func (m *mockSQLResult) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}
