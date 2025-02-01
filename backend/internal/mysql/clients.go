package mysql

import (
	"context"
	"database/sql"
	"strconv"
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
		Active:        true,
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

func (r *ClientRepository) UpdateClient(ctx context.Context, client *repository.UpdateClient) (error) {
	params := generated.UpdateClientParams{
		ID: client.ID,
		UpdatedBy: client.UpdatedBy,
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

	if client.Active != nil {
		params.Active = sql.NullBool{
			Valid: true,
			Bool: *client.Active,
		}
	}

	if client.BranchID != nil {
		params.BranchID = sql.NullInt32{
			Valid: true,
			Int32: int32(*client.BranchID),
		}
	}

	_, err := r.queries.UpdateClient(ctx, params)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update client: %s", err.Error())
	}

	return  nil
}

func (r *ClientRepository) UpdateClientOverpayment(ctx context.Context, phoneNumber string, overpayment float64) error {
	_, err := r.queries.UpdateClientOverpayment(ctx, generated.UpdateClientOverpaymentParams{
		PhoneNumber: phoneNumber,
		Overpayment: overpayment,
	})
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update client overpayment: %s", err.Error())
	}

	return nil
}

func (r *ClientRepository) ListClients(ctx context.Context, category *repository.ClientCategorySearch, pgData *pkg.PaginationMetadata) ([]repository.Client, pkg.PaginationMetadata, error) {
	params := generated.ListClientsByCategoryParams{
		Column1: "",
		FullName: "",
		PhoneNumber: "",
		Limit:    int32(pgData.PageSize),
		Offset:   int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
	}

	params2 := generated.CountClientsByCategoryParams{
		Column1: "",
		FullName: "",
		PhoneNumber: "",
	}

	if category.Search != nil {
		searchValue := "%" + *category.Search + "%"
		params.Column1 = "has_search"
		params.FullName = searchValue
		params.PhoneNumber = searchValue

		params2.Column1 = "has_search"
		params2.FullName = searchValue
		params2.PhoneNumber = searchValue
	}

	if category.Active != nil {
		params.Active = sql.NullBool{
			Valid: category.Active != nil,
			Bool: *category.Active,
		}
		params2.Active = sql.NullBool{
			Valid: category.Active != nil,
			Bool: *category.Active,
		}
	}

	clients, err := r.queries.ListClientsByCategory(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no clients found")
		}

		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list clients: %s", err.Error())
	}

	totalClients, err := r.queries.CountClientsByCategory(ctx, params2)
	if err != nil {
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get total loans: %s", err.Error())
	}

	result := make([]repository.Client, len(clients))
	for i, client := range clients {
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

		dueAmountBtye, _ := client.Dueamount.([]byte)
		dueAmount, _ := strconv.ParseFloat(string(dueAmountBtye), 64)

		result[i] = repository.Client{
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
			BranchName: &client.BranchName,
			DueAmount: dueAmount,
		}
	}

	return result, pkg.CreatePaginationMetadata(uint32(totalClients), pgData.PageSize, pgData.CurrentPage), nil
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

func (r *ClientRepository) GetClientIDByPhoneNumber(ctx context.Context, phoneNumber string) (uint32, error) {
	id, err := r.queries.GetClientIDByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, pkg.Errorf(pkg.NOT_FOUND_ERROR, "client not found")
		}

		return 0, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get client: %s", err.Error())
	}

	return id, nil
}

func (r *ClientRepository) ListClientsByBranch(ctx context.Context, branchID uint32, pgData *pkg.PaginationMetadata) ([]repository.Client, error) {
	clients, err := r.queries.ListClientsByBranch(ctx, generated.ListClientsByBranchParams{
		Limit:    int32(pgData.PageSize),
		Offset:   int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
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
		Limit:  int32(pgData.PageSize),
		Offset: int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
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
