package mysql

import (
	"context"
	"database/sql"
	"strings"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.NonPostedRepository = (*NonPostedRepository)(nil)

type NonPostedRepository struct {
	db      *Store
	queries generated.Querier
}

func NewNonPostedRepository(db *Store) *NonPostedRepository {
	return &NonPostedRepository{
		db:      db,
		queries: generated.New(db.db),
	}
}

func (r *NonPostedRepository) CreateNonPosted(ctx context.Context, nonPosted *repository.NonPosted) (repository.NonPosted, error) {
	params := generated.CreateNonPostedParams{
		TransactionSource: generated.NonPostedTransactionSource(nonPosted.TransactionSource),
		TransactionNumber: nonPosted.TransactionNumber,
		AccountNumber:     nonPosted.AccountNumber,
		PhoneNumber:       nonPosted.PhoneNumber,
		PayingName:        nonPosted.PayingName,
		Amount:            nonPosted.Amount,
		PaidDate:          nonPosted.PaidDate,
		AssignedBy: 	   nonPosted.AssignedBy,
	}

	if nonPosted.AssignedTo != nil {
		params.AssignTo = sql.NullInt32{
			Valid: true,
			Int32: int32(*nonPosted.AssignedTo),
		}
	}

	execResult, err := r.queries.CreateNonPosted(ctx, params)
	if err != nil {
		return repository.NonPosted{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create non posted: %s", err.Error())
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.NonPosted{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	nonPosted.ID = uint32(id)

	return *nonPosted, nil
}

func (r *NonPostedRepository) GetNonPosted(ctx context.Context, id uint32) (repository.NonPosted, error) {
	nonPosted, err := r.queries.GetNonPosted(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.NonPosted{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no non posted found")
		}

		return repository.NonPosted{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get non posted: %s", err.Error())
	}

	return convertGenerateNonPosted(nonPosted), nil
}

func (r *NonPostedRepository) ListNonPosted(ctx context.Context, category *repository.NonPostedCategory, pgData *pkg.PaginationMetadata) ([]repository.NonPosted, pkg.PaginationMetadata, error) {
	params := generated.ListNonPostedByCategoryParams{
		Column1: "",
		PayingName: "",
		AccountNumber: "",
		TransactionNumber: "",
		Column5: "",
		FINDINSET: "",
		Limit:    int32(pgData.PageSize),
		Offset:   int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
	}

	params2 := generated.CountNonPostedByCategoryParams{
		Column1: "",
		PayingName: "",
		AccountNumber: "",
		TransactionNumber: "",
		Column5: "",
		FINDINSET: "",
	}

	if category.Search != nil {
		searchValue := "%" + *category.Search + "%"
		params.Column1 = "has_search"
		params.PayingName = searchValue
		params.AccountNumber = searchValue
		params.TransactionNumber = searchValue

		params2.Column1 = "has_search"
		params2.PayingName = searchValue
		params2.AccountNumber = searchValue
		params2.TransactionNumber = searchValue
	}

	if category.Sources != nil {
		params.Column5 = "has_source"
		params2.Column5 = "has_source"		
		params.FINDINSET = *category.Sources
		params2.FINDINSET = *category.Sources
	}

	nonPosteds, err := r.queries.ListNonPostedByCategory(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no non posted found")
		}

		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list non posted: %s", err.Error())
	}

	totalNonPosted, err := r.queries.CountNonPostedByCategory(ctx, params2)
	if err != nil {
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get total loans: %s", err.Error())
	}

	rslt := make([]repository.NonPosted, len(nonPosteds))

	for i, nonPosted := range nonPosteds {
		rslt[i] = convertGenerateNonPosted(nonPosted)
	}

	return rslt, pkg.CreatePaginationMetadata(uint32(totalNonPosted), pgData.PageSize, pgData.CurrentPage), nil
}

func (r *NonPostedRepository) ListUnassignedNonPosted(ctx context.Context, pgData *pkg.PaginationMetadata) ([]repository.NonPosted, error) {
	nonPosteds, err := r.queries.ListUnassignedNonPosted(ctx, generated.ListUnassignedNonPostedParams{
		Limit:  int32(pgData.PageSize),
		Offset: int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no non posted found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list non posted: %s", err.Error())
	}

	rslt := make([]repository.NonPosted, len(nonPosteds))

	for i, nonPosted := range nonPosteds {
		rslt[i] = convertGenerateNonPosted(nonPosted)
	}

	return rslt, nil
}

func (r *NonPostedRepository) ListNonPostedByTransactionSource(
	ctx context.Context,
	transactionSource string,
	pgData *pkg.PaginationMetadata,
) ([]repository.NonPosted, error) {
	nonPosteds, err := r.queries.ListNonPostedByTransactionSource(ctx, generated.ListNonPostedByTransactionSourceParams{
		TransactionSource: generated.NonPostedTransactionSource(strings.ToUpper(transactionSource)),
		Limit:             int32(pgData.PageSize),
		Offset:            int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no non posted found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list non posted: %s", err.Error())
	}

	rslt := make([]repository.NonPosted, len(nonPosteds))

	for i, nonPosted := range nonPosteds {
		rslt[i] = convertGenerateNonPosted(nonPosted)
	}

	return rslt, nil
}

func (r *NonPostedRepository) DeleteNonPosted(ctx context.Context, id uint32) error {
	err := r.queries.DeleteNonPosted(ctx, id)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to delete non posted: %s", err.Error())
	}

	return nil
}

func (r *NonPostedRepository) GetReportPaymentData(ctx context.Context, filters services.ReportFilters) ([]services.PaymentReportData, services.PaymentSummary, error) {
	nonPosteds, err := r.queries.GetPaymentReportData(ctx, generated.GetPaymentReportDataParams{
		FromPaidDate: filters.StartDate,
		ToPaidDate:   filters.EndDate,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, services.PaymentSummary{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no non-posted payments found")
		}

		return nil, services.PaymentSummary{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get report payment data: %s", err.Error())
	}

	rslt := make([]services.PaymentReportData, len(nonPosteds))

	var totalPayments int64
	var totalAmountReceived float64
	sourceCount := make(map[string]int64)  // Count of each transaction source
	staffCount := make(map[string]int64)   // Count of each assigned staff

	for i, nonPosted := range nonPosteds {
		rslt[i] = services.PaymentReportData{
			TransactionSource: string(nonPosted.TransactionSource),
			TransactionNumber: nonPosted.TransactionNumber,
			AccountNumber:     nonPosted.AccountNumber,
			PayingName:        nonPosted.PayingName,
			Amount:            nonPosted.Amount,
			PaidDate:          nonPosted.PaidDate,
			AssignedTo:        nonPosted.AssignedName,
			AssignedBy:        nonPosted.AssignedBy,
		}

		totalPayments++
		totalAmountReceived += nonPosted.Amount
		sourceCount[string(nonPosted.TransactionSource)]++
		staffCount[nonPosted.AssignedBy]++
	}

	var mostCommonSource string
	var maxSourceCount int64
	for source, count := range sourceCount {
		if count > maxSourceCount {
			maxSourceCount = count
			mostCommonSource = source
		}
	}

	var mostAssignedStaff string
	var maxStaffCount int64
	for staff, count := range staffCount {
		if count > maxStaffCount {
			maxStaffCount = count
			mostAssignedStaff = staff
		}
	}

	summary := services.PaymentSummary{
		TotalPayments:       totalPayments,
		TotalAmountReceived: totalAmountReceived,
		MostCommonSource:    mostCommonSource,
		MostAssignedStaff:   mostAssignedStaff,
	}

	return rslt, summary, nil
}

func (r *NonPostedRepository) GetClientNonPosted(ctx context.Context, id uint32, phoneNumber string, pgData *pkg.PaginationMetadata) (repository.ClientNonPosted, pkg.PaginationMetadata, error) {
	var client generated.Client
	var loan generated.GetActiveLoanDetailsRow
	var installments []generated.Installment

	var clientPresent bool
	var clientHasNonPosted bool
	var clientHasActiveLoan bool
	var err error
	
	if id != 0 {
		client, err = r.queries.GetClient(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				clientPresent = false
			} else {
				return repository.ClientNonPosted{}, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get client: %s", err.Error())
			}
		}
		clientPresent = true
	}

	params := generated.GetClientsNonPostedParams{
		Limit:    int32(pgData.PageSize),
		Offset:   int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
	}

	params2 := generated.GetTotalPaidByIDorAccountNoParams{}
	params3 := generated.CountClientsNonPostedParams{}

	if clientPresent {
		params.AssignTo = sql.NullInt32{
			Valid: true,
			Int32: int32(client.ID),
		}
		params2.AssignTo = sql.NullInt32{
			Valid: true,
			Int32: int32(client.ID),
		}
		params3.AssignTo = sql.NullInt32{
			Valid: true,
			Int32: int32(client.ID),
		}
	} else {
		if phoneNumber == "" {
			return repository.ClientNonPosted{}, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INVALID_ERROR, "both id and phonenumber cannot be missing")
		}
		params.AccountNumber = sql.NullString{
			Valid: true,
			String: phoneNumber,
		}
		params2.AccountNumber = sql.NullString{
			Valid: true,
			String: phoneNumber,
		}
		params3.AccountNumber = sql.NullString{
			Valid: true,
			String: phoneNumber,
		}
	}

	nonPosteds, err := r.queries.GetClientsNonPosted(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			clientHasNonPosted = false
		} else {
			return repository.ClientNonPosted{}, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get non-posted: %s", err.Error())
		}
	}
	clientHasNonPosted = true

	totalNonPosted, err := r.queries.CountClientsNonPosted(ctx, params3)
	if err != nil {
		return repository.ClientNonPosted{}, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get total loans: %s", err)
	}

	totalPaid, err := r.queries.GetTotalPaidByIDorAccountNo(ctx, params2)
	if err != nil {
		return repository.ClientNonPosted{}, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get total payment: %s", err)
	}

	if clientPresent {
		loan, err = r.queries.GetActiveLoanDetails(ctx, client.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				clientHasActiveLoan = false
			} else {
				return repository.ClientNonPosted{}, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loan: %s", err)
			}
		}
		clientHasActiveLoan = true
	}

	if clientHasActiveLoan {
		installments, err = r.queries.ListInstallmentsByLoan(ctx, loan.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				return repository.ClientNonPosted{}, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "loan has no  installemt!!!")
			} else {
				return repository.ClientNonPosted{}, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loan installments: %s", err.Error())
			}
		}
	}

	rslt := repository.ClientNonPosted{}
	if clientPresent {
		branch, _ := r.queries.GetBranch(ctx, client.BranchID)
		rslt.ClientDetails.ID = client.ID
		rslt.ClientDetails.FullName = client.FullName
		rslt.ClientDetails.PhoneNumber = client.PhoneNumber
		rslt.ClientDetails.Overpayment = client.Overpayment
		rslt.ClientDetails.BranchName = branch.Name
	}

	if clientHasActiveLoan {
		rslt.LoanDetails.ID = loan.ID
		rslt.LoanDetails.LoanAmount = loan.LoanAmount
		rslt.LoanDetails.RepayAmount = loan.RepayAmount
		rslt.LoanDetails.DisbursedOn = loan.DisbursedOn.Time.Format("2006-02-01")
		rslt.LoanDetails.DueDate = loan.DueDate.Time.Format("2006-02-01")
		rslt.LoanDetails.PaidAmount = loan.PaidAmount
		rslt.LoanDetails.Installments = convertGeneratedInstallmentList(installments)
	}

	if clientHasNonPosted {
		paymentDetails := make([]repository.NonPostedShort, len(nonPosteds))
		for i, nonPosted := range nonPosteds {
			paymentDetails[i] = repository.NonPostedShort{
				ID:                nonPosted.ID,
				TransactionSource: string(nonPosted.TransactionSource),
				TransactionNumber: nonPosted.TransactionNumber,
				AccountNumber:     nonPosted.AccountNumber,
				PhoneNumber:       nonPosted.PhoneNumber,
				PayingName:        nonPosted.PayingName,
				Amount:            nonPosted.Amount,
				PaidDate:          nonPosted.PaidDate,
				AssignedBy: 	   nonPosted.AssignedBy,
			}
		}
		rslt.PaymentDetails = paymentDetails
		rslt.TotalPaid = pkg.InterfaceFloat64(totalPaid)
	}

	return rslt, pkg.CreatePaginationMetadata(uint32(totalNonPosted), pgData.PageSize, pgData.CurrentPage), nil
}

func convertGenerateNonPosted(nonPosted generated.NonPosted) repository.NonPosted {
	var assignedTo *uint32

	if nonPosted.AssignTo.Valid {
		value := uint32(nonPosted.AssignTo.Int32)
		assignedTo = &value
	}

	return repository.NonPosted{
		ID:                nonPosted.ID,
		TransactionSource: string(nonPosted.TransactionSource),
		TransactionNumber: nonPosted.TransactionNumber,
		AccountNumber:     nonPosted.AccountNumber,
		PhoneNumber:       nonPosted.PhoneNumber,
		PayingName:        nonPosted.PayingName,
		Amount:            nonPosted.Amount,
		PaidDate:          nonPosted.PaidDate,
		AssignedTo:        assignedTo,
	}
}
