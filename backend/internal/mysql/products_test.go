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
	store := NewStore()

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
		wantResult []*repository.Product
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
			wantResult: []*repository.Product{{ID: 42, BranchID: 1, LoanAmount: 1000}},
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
