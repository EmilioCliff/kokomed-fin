package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.ClientRepository = (*ClientRepository)(nil)

type ClientRepository struct {
	db      *Store
	queries generated.Querier
}

func NewClientRepository(db *Store) *ClientRepository {
	return &ClientRepository{
		db:      db,
		queries: generated.New(db.db),
	}
}

func (r *ClientRepository) CreateClient(ctx context.Context, client *repository.Client) (repository.Client, error) {
	params := generated.CreateClientParams{
		FullName:      client.FullName,
		PhoneNumber:   client.PhoneNumber,
		Gender:        generated.ClientsGender(client.Gender),
		Active:        client.Active,
		BranchID:      client.BranchID,
		AssignedStaff: client.AssignedStaff,
		UpdatedBy:     client.UpdatedBy,
		CreatedBy:     client.CreatedBy,
	}

	if client.IdNumber != nil {
		params.IDNumber = sql.NullString{
			String: *client.IdNumber,
			Valid:  true,
		}
	}

	if client.Dob != nil {
		params.Dob = sql.NullTime{
			Time:  *client.Dob,
			Valid: true,
		}
	}

	execResult, err := r.queries.CreateClient(ctx, params)
	if err != nil {
		return repository.Client{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create client: %s", err.Error())
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.Client{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	client.ID = uint32(id)

	return *client, nil
}

func (r *ClientRepository) UpdateClient(ctx context.Context, client *repository.Client) (repository.Client, error) {
	params := generated.UpdateClientParams{
		FullName:      client.FullName,
		PhoneNumber:   client.PhoneNumber,
		Gender:        generated.ClientsGender(client.Gender),
		Active:        client.Active,
		BranchID:      client.BranchID,
		AssignedStaff: client.AssignedStaff,
		UpdatedBy:     client.UpdatedBy,
	}

	if client.IdNumber != nil {
		params.IDNumber = sql.NullString{
			String: *client.IdNumber,
			Valid:  true,
		}
	}

	if client.Dob != nil {
		params.Dob = sql.NullTime{
			Time:  *client.Dob,
			Valid: true,
		}
	}

	execResult, err := r.queries.UpdateClient(ctx, params)
	if err != nil {
		return repository.Client{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update client: %s", err.Error())
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.Client{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	client.ID = uint32(id)

	return *client, nil
}

func (r *ClientRepository) ListClients(ctx context.Context, pgData *pkg.PaginationMetadata) ([]repository.Client, error) {
	clients, err := r.queries.ListClients(ctx, generated.ListClientsParams{
		Limit:  pkg.GetPageSize(),
		Offset: pkg.CalculateOffset(pgData.CurrentPage),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no clients found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list clients: %s", err.Error())
	}

	result := make([]repository.Client, len(clients))
	for i, client := range clients {
		result[i] = convertGeneratedClient(client)
	}

	return result, nil
}

func (r *ClientRepository) GetClient(ctx context.Context, clientID uint32) (repository.Client, error) {
	client, err := r.queries.GetClient(ctx, clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.Client{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "client not found")
		}

		return repository.Client{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get client: %s", err.Error())
	}

	return convertGeneratedClient(client), nil
}

func (r *ClientRepository) GetClientByPhoneNumber(ctx context.Context, phoneNumber string) (repository.Client, error) {
	client, err := r.queries.GetClientByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.Client{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "client not found")
		}

		return repository.Client{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get client: %s", err.Error())
	}

	return convertGeneratedClient(client), nil
}

func (r *ClientRepository) ListClientsByBranch(ctx context.Context, branchID uint32, pgData *pkg.PaginationMetadata) ([]repository.Client, error) {
	clients, err := r.queries.ListClientsByBranch(ctx, generated.ListClientsByBranchParams{
		Limit:    pkg.GetPageSize(),
		Offset:   pkg.CalculateOffset(pgData.CurrentPage),
		BranchID: branchID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no clients found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list clients: %s", err.Error())
	}

	result := make([]repository.Client, len(clients))
	for i, client := range clients {
		result[i] = convertGeneratedClient(client)
	}

	return result, nil
}

func (r *ClientRepository) ListClientsByActiveStatus(ctx context.Context, active bool, pgData *pkg.PaginationMetadata) ([]repository.Client, error) {
	clients, err := r.queries.ListClientsByActiveStatus(ctx, generated.ListClientsByActiveStatusParams{
		Limit:  pkg.GetPageSize(),
		Offset: pkg.CalculateOffset(pgData.CurrentPage),
		Active: active,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no clients found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list clients: %s", err.Error())
	}

	result := make([]repository.Client, len(clients))
	for i, client := range clients {
		result[i] = convertGeneratedClient(client)
	}

	return result, nil
}

func convertGeneratedClient(client generated.Client) repository.Client {
	var dob *time.Time

	if client.Dob.Valid {
		value := client.Dob.Time
		dob = &value
	}

	var idNo *string

	if client.IDNumber.Valid {
		value := client.IDNumber.String
		idNo = &value
	}

	return repository.Client{
		ID:            client.ID,
		FullName:      client.FullName,
		PhoneNumber:   client.PhoneNumber,
		IdNumber:      idNo,
		Dob:           dob,
		Gender:        string(client.Gender),
		Active:        client.Active,
		BranchID:      client.BranchID,
		AssignedStaff: client.AssignedStaff,
		Overpayment:   client.Overpayment,
		UpdatedBy:     client.UpdatedBy,
		UpdatedAt:     client.UpdatedAt,
		CreatedBy:     client.CreatedBy,
		CreatedAt:     client.CreatedAt,
	}
}
