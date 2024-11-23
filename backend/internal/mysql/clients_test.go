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

func NewTestClientRepository() *ClientRepository {
	store := NewStore(pkg.Config{})

	return NewClientRepository(store)
}

func TestClientRepository_CreateClient(t *testing.T) {
	r := NewTestClientRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.Client
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					CreateClient(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&mockSQLResult{lastInsertID: 1, rowsAffected: 1}, nil)
			},
			wantErr:    false,
			err:        nil,
			wantResult: repository.Client{ID: 1, FullName: "test"},
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					CreateClient(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&mockSQLResult{lastInsertID: 1, rowsAffected: 1}, errors.New("error"))
			},
			wantErr:    true,
			err:        errors.New(pkg.INTERNAL_ERROR),
			wantResult: repository.Client{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.CreateClient(context.Background(), &repository.Client{
				FullName: "test",
			})

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

func TestClientRepository_UpdateClient(t *testing.T) {
	r := NewTestClientRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.Client
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					UpdateClient(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&mockSQLResult{lastInsertID: 1, rowsAffected: 1}, nil)
			},
			wantErr:    false,
			err:        nil,
			wantResult: repository.Client{ID: 1, FullName: "test"},
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					UpdateClient(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&mockSQLResult{lastInsertID: 1, rowsAffected: 1}, errors.New("error"))
			},
			wantErr:    true,
			err:        errors.New(pkg.INTERNAL_ERROR),
			wantResult: repository.Client{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.UpdateClient(context.Background(), &repository.Client{
				ID:       1,
				FullName: "test",
			})

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

func TestClientRepository_ListClients(t *testing.T) {
	r := NewTestClientRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult []repository.Client
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListClients(gomock.Any(), gomock.AssignableToTypeOf(generated.ListClientsParams{})).
					Times(1).
					Return([]generated.Client{
						{
							ID:       1,
							FullName: "test",
						},
					}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: []repository.Client{
				{
					ID:       1,
					FullName: "test",
				},
			},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListClients(gomock.Any(), gomock.AssignableToTypeOf(generated.ListClientsParams{})).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			wantErr:    true,
			err:        errors.New(pkg.NOT_FOUND_ERROR),
			wantResult: nil,
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListClients(gomock.Any(), gomock.AssignableToTypeOf(generated.ListClientsParams{})).
					Times(1).
					Return(nil, errors.New("error"))
			},
			wantErr:    true,
			err:        errors.New(pkg.INTERNAL_ERROR),
			wantResult: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.ListClients(context.Background(), &pkg.PaginationMetadata{})

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

func TestClientRepository_GetClient(t *testing.T) {
	r := NewTestClientRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.Client
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetClient(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.Client{ID: 1, FullName: "test"}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: repository.Client{
				ID:       1,
				FullName: "test",
			},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetClient(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.Client{}, sql.ErrNoRows)
			},
			wantErr:    true,
			err:        errors.New(pkg.NOT_FOUND_ERROR),
			wantResult: repository.Client{},
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetClient(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.Client{}, errors.New("error"))
			},
			wantErr:    true,
			err:        errors.New(pkg.INTERNAL_ERROR),
			wantResult: repository.Client{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.GetClient(context.Background(), 1)

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

func TestClientRepository_GetClientByPhoneNumber(t *testing.T) {
	r := NewTestClientRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult repository.Client
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetClientByPhoneNumber(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.Client{ID: 1, FullName: "test"}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: repository.Client{
				ID:       1,
				FullName: "test",
			},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetClientByPhoneNumber(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.Client{}, sql.ErrNoRows)
			},
			wantErr:    true,
			err:        errors.New(pkg.NOT_FOUND_ERROR),
			wantResult: repository.Client{},
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					GetClientByPhoneNumber(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generated.Client{}, errors.New("error"))
			},
			wantErr:    true,
			err:        errors.New(pkg.INTERNAL_ERROR),
			wantResult: repository.Client{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.GetClientByPhoneNumber(context.Background(), "test")

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

func TestClientRepository_ListClientsByBranch(t *testing.T) {
	r := NewTestClientRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult []repository.Client
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListClientsByBranch(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.Client{
						{ID: 1, FullName: "test"},
					}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: []repository.Client{
				{ID: 1, FullName: "test"},
			},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListClientsByBranch(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.Client{}, sql.ErrNoRows)
			},
			wantErr:    true,
			err:        errors.New(pkg.NOT_FOUND_ERROR),
			wantResult: []repository.Client{},
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListClientsByBranch(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.Client{}, errors.New("error"))
			},
			wantErr:    true,
			err:        errors.New(pkg.INTERNAL_ERROR),
			wantResult: []repository.Client{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.ListClientsByBranch(context.Background(), 1, &pkg.PaginationMetadata{})

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

func TestClientRepository_ListClientsByActiveStatus(t *testing.T) {
	r := NewTestClientRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)
	r.queries = mockQueries

	tests := []struct {
		name       string
		buildStubs func(mockQueries *mockdb.MockQuerier)
		wantErr    bool
		err        error
		wantResult []repository.Client
	}{
		{
			name: "OK",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListClientsByActiveStatus(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.Client{
						{ID: 1, FullName: "test"},
					}, nil)
			},
			wantErr: false,
			err:     nil,
			wantResult: []repository.Client{
				{ID: 1, FullName: "test"},
			},
		},
		{
			name: "Not Found",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListClientsByActiveStatus(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.Client{}, sql.ErrNoRows)
			},
			wantErr:    true,
			err:        errors.New(pkg.NOT_FOUND_ERROR),
			wantResult: []repository.Client{},
		},
		{
			name: "Internal Error",
			buildStubs: func(mockQueries *mockdb.MockQuerier) {
				mockQueries.EXPECT().
					ListClientsByActiveStatus(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]generated.Client{}, errors.New("error"))
			},
			wantErr:    true,
			err:        errors.New(pkg.INTERNAL_ERROR),
			wantResult: []repository.Client{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries)

			result, err := r.ListClientsByActiveStatus(context.Background(), true, &pkg.PaginationMetadata{})

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
