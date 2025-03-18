package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"go.opentelemetry.io/otel/codes"
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

func (r *ClientRepository) CreateClient(ctx context.Context, client *repository.Client) (repository.ClientFullData, error) {
	tc, span := r.db.tracer.Start(ctx, "Client Repo: CreateClient")
	defer span.End()
	
	params := generated.CreateClientParams{
		FullName:      client.FullName,
		PhoneNumber:   client.PhoneNumber,
		Gender:        generated.ClientsGender(client.Gender),
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

	execResult, err := r.queries.CreateClient(tc, params)
	if err != nil {
		setSpanError(span, codes.Error, err, "failed to create client")
		return repository.ClientFullData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create client: %s", err.Error())
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.ClientFullData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	client.ID = uint32(id)

	return r.GetClientFullData(tc, client.ID)
}

func (r *ClientRepository) GetClientFullData(ctx context.Context, id uint32) (repository.ClientFullData, error) {
	tc, span := r.db.tracer.Start(ctx, "Client Repo: GetClientFullData")
	defer span.End()

	client, err := r.queries.GetClientFullData(tc, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.ClientFullData{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "client not found")
		}

		setSpanError(span, codes.Error, err, "failed to get client full data")
		return repository.ClientFullData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get client full data: %s", err.Error())
	}

	return convertClientFullData(client), nil
}

func (r *ClientRepository) UpdateClient(ctx context.Context, client *repository.UpdateClient) (error) {
	tc, span := r.db.tracer.Start(ctx, "Client Repo: UpdateClient")
	defer span.End()

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

	_, err := r.queries.UpdateClient(tc, params)
	if err != nil {
		setSpanError(span, codes.Error, err, "failed to update client")
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update client: %s", err.Error())
	}

	return  nil
}

// used during transactions
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

func (r *ClientRepository) ListClients(ctx context.Context, category *repository.ClientCategorySearch, pgData *pkg.PaginationMetadata) ([]repository.ClientFullData, pkg.PaginationMetadata, error) {
	tc, span := r.db.tracer.Start(ctx, "Client Repo: ListClients")
	defer span.End()

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

	clients, err := r.queries.ListClientsByCategory(tc, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no clients found")
		}

		setSpanError(span, codes.Error, err, "failed to list clients")
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list clients: %s", err.Error())
	}

	totalClients, err := r.queries.CountClientsByCategory(ctx, params2)
	if err != nil {
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get total loans: %s", err.Error())
	}

	result := make([]repository.ClientFullData, len(clients))
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

		result[i] = repository.ClientFullData{
			ID:            client.ID,
			FullName:      client.FullName,
			PhoneNumber:   client.PhoneNumber,
			IDNumber:      idNo,
			DOB:           dob,
			Gender:        string(client.Gender),
			Active:        client.Active,
			AssignedStaff: repository.UserShortResponse{
				ID: uint32(client.AssignedUserID.Int32),
				FullName: client.AssignedUserName.String,
				PhoneNumber: client.AssignedUserPhone.String,
				Email: client.AssignedUserEmail.String,
				Role: string(client.AssignedUserRole.UsersRole),
			},
			Overpayment:   client.Overpayment,
			UpdatedBy:     repository.UserShortResponse{
				ID: uint32(client.UpdatedUserID.Int32),
				FullName: client.UpdatedUserName.String,
				PhoneNumber: client.UpdatedUserPhone.String,
				Email: client.UpdatedUserEmail.String,
				Role: string(client.UpdatedUserRole.UsersRole),
			},
			UpdatedAt:     client.UpdatedAt,
			CreatedBy:     repository.UserShortResponse{
				ID: uint32(client.CreatedUserID.Int32),
				FullName: client.CreatedUserName.String,
				PhoneNumber: client.CreatedUserPhone.String,
				Email: client.CreatedUserEmail.String,
				Role: string(client.CreatedUserRole.UsersRole),
			},
			CreatedAt:     client.CreatedAt,
			BranchName:    client.BranchName,
			DueAmount: dueAmount,
		}
	}

	return result, pkg.CreatePaginationMetadata(uint32(totalClients), pgData.PageSize, pgData.CurrentPage), nil
}

// func (r *ClientRepository) GetClient(ctx context.Context, clientID uint32) (repository.ClientFullData, error) {
// 	tc, span := r.db.tracer.Start(ctx, "Client Repo: GetClient")
// 	defer span.End()

// 	client, err := r.queries.GetClient(tc, clientID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return repository.ClientFullData{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "client not found")
// 		}

// 		setSpanError(span, codes.Error, err, "failed to get client")
// 		return repository.ClientFullData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get client: %s", err.Error())
// 	}

// 	return convertGeneratedClient(client), nil
// }

func (r *ClientRepository) GetClientIDByPhoneNumber(ctx context.Context, phoneNumber string) (uint32, error) {
	tc, span := r.db.tracer.Start(ctx, "Client Repo: GetClientIDByPhoneNumber")
	defer span.End()

	id, err := r.queries.GetClientIDByPhoneNumber(tc, phoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, pkg.Errorf(pkg.NOT_FOUND_ERROR, "client not found")
		}

		setSpanError(span, codes.Error, err, "failed to get client")
		return 0, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get client: %s", err.Error())
	}

	return id, nil
}

func (r *ClientRepository) ListClientsByBranch(ctx context.Context, branchID uint32, pgData *pkg.PaginationMetadata) ([]repository.Client, error) {
	tc, span := r.db.tracer.Start(ctx, "Client Repo: ListClientsByBranch")
	defer span.End()

	clients, err := r.queries.ListClientsByBranch(tc, generated.ListClientsByBranchParams{
		Limit:    int32(pgData.PageSize),
		Offset:   int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
		BranchID: branchID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no clients found")
		}

		setSpanError(span, codes.Error, err, "failed to list clients")
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list clients: %s", err.Error())
	}

	result := make([]repository.Client, len(clients))
	for i, client := range clients {
		result[i] = convertGeneratedClient(client)
	}

	return result, nil
}

func (r *ClientRepository) ListClientsByActiveStatus(ctx context.Context, active bool, pgData *pkg.PaginationMetadata) ([]repository.Client, error) {
	tc, span := r.db.tracer.Start(ctx, "Client Repo: ListClientsByActiveStatus")
	defer span.End()

	clients, err := r.queries.ListClientsByActiveStatus(tc, generated.ListClientsByActiveStatusParams{
		Limit:  int32(pgData.PageSize),
		Offset: int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
		Active: active,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no clients found")
		}

		setSpanError(span, codes.Error, err, "failed to list clients")
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list clients: %s", err.Error())
	}

	result := make([]repository.Client, len(clients))
	for i, client := range clients {
		result[i] = convertGeneratedClient(client)
	}

	return result, nil
}

func (r *ClientRepository) GetReportClientAdminData(ctx context.Context, filters services.ReportFilters) ([]services.ClientAdminsReportData, services.ClientSummary, error) {
	tc, span := r.db.tracer.Start(ctx, "Client Repo: GetReportClientAdminData")
	defer span.End()

	clients, err := r.GetClientAdminsReportData(tc, GetClientAdminsReportDataParams{
		StartDate: filters.StartDate,
		EndDate:   filters.EndDate,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, services.ClientSummary{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no client found")
		}

		setSpanError(span, codes.Error, err, "failed to get report client admin data")
		return nil, services.ClientSummary{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get report client admin data: %s", err.Error())
	}

	rslt := make([]services.ClientAdminsReportData, len(clients))

	var totalClients int64
	var totalDisbursed, totalPaid, totalOwed float64
	var mostLoansClient, highestPayingClient string
	var maxLoansGiven, maxPaidAmount float64

	for i, client := range clients {
		totalPaidAmount := pkg.InterfaceFloat64(client.TotalPaid)
		totalDisbursedAmount := pkg.InterfaceFloat64(client.TotalDisbursed)
		totalOwedAmount := pkg.InterfaceFloat64(client.TotalOwed)

		rslt[i] = services.ClientAdminsReportData{
			Name:            client.Name,
			BranchName:      client.BranchName.String,
			TotalLoanGiven:  client.TotalLoanGiven,
			DefaultedLoans:  client.DefaultedLoans,
			ActiveLoans:     client.ActiveLoans,
			CompletedLoans:  client.CompletedLoans,
			InactiveLoans:   client.InactiveLoans,
			Overpayment:     client.Overpayment,
			PhoneNumber:     client.PhoneNumber,
			TotalPaid:       totalPaidAmount,
			TotalDisbursed:  totalDisbursedAmount,
			TotalOwed:       totalOwedAmount,
			RateScore:       pkg.InterfaceFloat64(client.RateScore),
			DefaultRate:     pkg.InterfaceFloat64(client.DefaultRate),
		}
		totalClients++
		totalDisbursed += totalDisbursedAmount
		totalPaid += totalPaidAmount
		totalOwed += totalOwedAmount

		if float64(client.TotalLoanGiven) > maxLoansGiven {
			maxLoansGiven = float64(client.TotalLoanGiven)
			mostLoansClient = client.Name
		}

		if totalPaidAmount > maxPaidAmount {
			maxPaidAmount = totalPaidAmount
			highestPayingClient = client.Name
		}
	}

	summary := services.ClientSummary{
		TotalClients:        totalClients,
		MostLoansClient:     mostLoansClient,
		HighestPayingClient: highestPayingClient,
		TotalDisbursed:      totalDisbursed,
		TotalPaid:           totalPaid,
		TotalOwed:           totalOwed,
	}

	return rslt, summary, nil
}

func (r *ClientRepository) GetReportClientClientsData(ctx context.Context,id uint32, filters services.ReportFilters) (services.ClientClientsReportData, error) {
	tc, span := r.db.tracer.Start(ctx, "Client Repo: GetReportClientClientsData")
	defer span.End()

	client, err := r.GetClientClientsReportData(tc, GetClientClientsReportDataParams{
		StartDate: filters.StartDate,
		EndDate: filters.EndDate,
		ID: id,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return services.ClientClientsReportData{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no client found")
		}

		setSpanError(span, codes.Error, err, "failed to get report client admin data")
		return services.ClientClientsReportData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get report client admin data: %s", err.Error())
	}

	return convertClientReportData(client)
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

func convertClientFullData(client generated.GetClientFullDataRow) repository.ClientFullData {
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

	return repository.ClientFullData{
		ID:			client.ClientID,
		FullName:			client.ClientName,
		PhoneNumber:			client.ClientPhone,
		IDNumber:			idNo,
		DOB:			dob,
		Gender:			string(client.Gender),
		Active:			client.Active,
		BranchName:			client.BranchName,
		Overpayment:			client.Overpayment,
		CreatedAt:			client.ClientCreatedAt,
		AssignedStaff: repository.UserShortResponse{
			ID: client.AssignedUserID,
			FullName: client.AssignedUserName,
			PhoneNumber: client.AssignedUserPhone,
			Email: client.AssignedUserEmail,
			Role: string(client.AssignedUserRole),
		},
		CreatedBy: repository.UserShortResponse{
			ID: client.CreatedByID,
			FullName: client.CreatedByName,
			PhoneNumber: client.CreatedByPhone,
			Email: client.CreatedByEmail,
			Role: string(client.CreatedByRole),
		},
	}
}

func convertClientReportData(row GetClientClientsReportDataRow) (services.ClientClientsReportData, error) {
	var loans []services.ClientClientReportDataLoans
	if row.Loans != nil {
		loansByte, ok := row.Loans.([]byte)
		if !ok {
			return services.ClientClientsReportData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to convert loans to bytes")
		}

		err := json.Unmarshal(loansByte, &loans)
		if err != nil {
			return services.ClientClientsReportData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "error unmarshalling loans: %v", err)
		}
	}

	var payments []services.ClientClientReportDataPayments
	if row.Payments != nil {
		paymentsByte, ok := row.Payments.([]byte)
		if !ok {
			return services.ClientClientsReportData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to convert payments to bytes")
		}

		err := json.Unmarshal(paymentsByte, &payments) 
		if err != nil {
			return services.ClientClientsReportData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "error unmarshalling payments: ", err)
		}
	}


	rsp := services.ClientClientsReportData{
		Name:          row.Name,
		PhoneNumber:   row.PhoneNumber,
		IDNumber:      nil,
		Dob:           nil,
		BranchName:    row.BranchName.String,
		AssignedStaff: row.AssignedStaff.String,
		Active:        "active",
		Loans:         loans,
		Payments:      payments,
	}

	if !row.Active { rsp.Active = "inactive"}
	if row.Dob.Valid {
		rsp.Dob = &row.Dob.Time
	}
	if row.IDNumber.Valid {
		rsp.IDNumber = &row.IDNumber.String
	}

	return rsp, nil
}