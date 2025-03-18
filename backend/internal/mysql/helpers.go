package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"go.opentelemetry.io/otel/codes"
)


var _ repository.HelperRepository = (*HelperRepository)(nil) 

type HelperRepository struct {
	db *Store
	queries generated.Querier
}

func NewHelperRepository(db *Store) *HelperRepository {
	return &HelperRepository{
		db:      db,
		queries: generated.New(db.db),
	}
}

func (r *HelperRepository) GetDashboardData(ctx context.Context) (repository.DashboardData, error) {
	tc, span := r.db.tracer.Start(ctx, "Helper Repo: GetDashboardData")
	defer span.End()

	widgetsData, err := r.queries.DashBoardDataHelper(tc)
	if err != nil {
		setSpanError(span, codes.Error, err, "failed to get dashboard data")
		return repository.DashboardData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting dashboard widgets")
	}

	loansData, err := r.queries.DashBoardInactiveLoans(tc)
	if err != nil {
		setSpanError(span, codes.Error, err, "failed to get dashboard inactive loans")
		return repository.DashboardData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting dashboard inactive loans")
	}

	paymentsData, err := r.queries.DashBoardRecentsPayments(tc)
	if err != nil {
		setSpanError(span, codes.Error, err, "failed to get dashboard recent payment")
		return repository.DashboardData{}, pkg.Errorf(pkg.INTERNAL_ERROR, err.Error())
	}

	inactiveLoans := make([]repository.InactiveLoan, len(loansData))
	for idx, loans := range loansData {
		inactiveLoans[idx] = repository.InactiveLoan{
			ID:          loans.ID,
			Amount:      loans.LoanAmount,
			ClientName: loans.ClientName,
			ApprovedByName: loans.ApprovedByName,
			RepayAmount: loans.RepayAmount,
			ApprovedOn:  loans.CreatedAt,
		}
	}

	recentPayments := make([]repository.Payment, len(paymentsData))
	for idx, payment := range paymentsData {
		recentPayments[idx] = repository.Payment{
			ID: payment.ID,
			PayingName: payment.PayingName,
			Amount: payment.Amount,
			PaidDate: payment.PaidDate,
		}
	}

		totalPaymentsReceived  := pkg.InterfaceFloat64(widgetsData.TotalPaymentsReceived)
		totalNonPosted := pkg.InterfaceFloat64(widgetsData.TotalNonPosted)

		widgets := []repository.Widget{
			{
				Title:"Customers",
				MainAmount: float64(widgetsData.TotalClients),
				Active: float64(widgetsData.ActiveClients),
				ActiveTitle: "Active",
				Closed: float64(widgetsData.TotalClients) - float64(widgetsData.ActiveClients),
				ClosedTitle: "Inactive",
			},
			{
				Title: "Loans",
				MainAmount: float64(widgetsData.TotalLoans),
				Active: float64(widgetsData.ActiveLoans),
				ActiveTitle: "Active",
				Closed: float64(widgetsData.InactiveLoans),
				ClosedTitle: "Inactive",
			},
			{
				Title: "Transactions",
				MainAmount: pkg.InterfaceFloat64(widgetsData.TotalLoanAmount),
				Active: pkg.InterfaceFloat64(widgetsData.TotalLoanDisbursed),
				ActiveTitle: "Disbursed",
				Closed: pkg.InterfaceFloat64(widgetsData.TotalLoanPaid),
				ClosedTitle: "Completed Loans",
				Currency: "Ksh",
			},
			{
				Title: "Payments",
				MainAmount: totalPaymentsReceived,
				Active: totalPaymentsReceived - totalNonPosted,
				ActiveTitle: "Posted",
				Closed: totalNonPosted,
				ClosedTitle: "Non-Posted",
				Currency: "Ksh",
			},
		}
		rsp := repository.DashboardData{
			WidgetData: widgets,
			InactiveLoans: inactiveLoans,
			RecentPayments: recentPayments,
		}
	return rsp, nil
}

func (r *HelperRepository) GetProductData(ctx context.Context) ([]repository.ProductData, error) {
	tc, span := r.db.tracer.Start(ctx, "Helper Repo: GetProductData")
	defer span.End()

	products, err := r.queries.HelperProduct(tc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		setSpanError(span, codes.Error, err, "failed to get product ddata")
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting product data")		
	}

	rsp := make([]repository.ProductData, len(products))
	for idx, product := range products {
		rsp[idx] = repository.ProductData{
			ID: product.Productid,
			Name: fmt.Sprintf("%.2f %s", product.Loanamount, product.Branchname),
		}
	}

	return rsp, nil
}
func (r *HelperRepository) GetClientData(ctx context.Context) ([]repository.ClientData, error) {
	tc, span := r.db.tracer.Start(ctx, "Helper Repo: GetClientData")
	defer span.End()

	clients, err := r.queries.HelperClient(tc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		setSpanError(span, codes.Error, err, "failed to get clients data")
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting clients list data")		
	}

	rsp := make([]repository.ClientData, len(clients))
	for idx, client := range clients {
		rsp[idx] = repository.ClientData{
			ID: client.ID,
			Name: fmt.Sprintf("%s - %s", client.FullName, client.PhoneNumber),
		}
	}

	return rsp, nil
}
func (r *HelperRepository) GetLoanOfficerData(ctx context.Context) ([]repository.LoanOfficerData, error) {
	tc, span := r.db.tracer.Start(ctx, "Helper Repo: GetLoanOfficerData")
	defer span.End()

	users, err := r.queries.HelperUser(tc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		setSpanError(span, codes.Error, err, "failed to get user data")
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting user list data")		
	}

	rsp := make([]repository.LoanOfficerData, len(users))
	for idx, user := range users {
		rsp[idx] = repository.LoanOfficerData{
			ID: user.ID,
			Name: user.FullName,
		}
	}

	return rsp, nil
}

func (r *HelperRepository) GetBranchData(ctx context.Context) ([]repository.BranchData,  error) {
	tc, span := r.db.tracer.Start(ctx, "Helper Repo: GetBranchData")
	defer span.End()

	branches, err := r.queries.ListBranches(tc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		setSpanError(span, codes.Error, err, "failed to get branch data")
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting branch list data")		
	}

	rsp := make([]repository.BranchData, len(branches))
	for idx, branch := range branches {
		rsp[idx] = repository.BranchData{
			ID: branch.ID,
			Name: branch.Name,
		}
	}

	return rsp, nil
}

func (r *HelperRepository) GetLoanData(ctx context.Context) ([]repository.LoanData,  error) {
	tc, span := r.db.tracer.Start(ctx, "Helper Repo: GetLoanData")
	defer span.End()

	loansId, err := r.queries.GetLoanData(tc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		setSpanError(span, codes.Error, err, "failed to get loan data")
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting loan list data")		
	}

	rsp := make([]repository.LoanData, len(loansId))
	for idx, id := range loansId {
		rsp[idx] = repository.LoanData{
			ID: id,
			Name: fmt.Sprintf("LN%03d", id),
		}
	}

	return rsp, nil
}

func (r *HelperRepository) GetUserFullname(ctx context.Context, id uint32) (string, error) {
	name, err := r.queries.HelperUserById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows{
			return "", pkg.Errorf(pkg.NOT_FOUND_ERROR, "user not found")
		}

		return "", pkg.Errorf(pkg.INTERNAL_ERROR, "error getting user fullname", err)
	}

	return name, err
}

func (r *HelperRepository) GetClientNonPayments(ctx context.Context,id uint32, phoneNumber string, pgData *pkg.PaginationMetadata) ([]repository.NonPostedShort, pkg.PaginationMetadata, error) {
	tc, span := r.db.tracer.Start(ctx, "Helper Repo: GetClientNonPayments")
	defer span.End()

	params := generated.GetClientsNonPostedParams{
		Limit:    int32(pgData.PageSize),
		Offset:   int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
	}

	params2 := generated.CountClientsNonPostedParams{}

	if id != 0 {
		params.AssignTo = sql.NullInt32{
			Valid: true,
			Int32: int32(id),
		}
		params2.AssignTo = sql.NullInt32{
			Valid: true,
			Int32: int32(id),
		}
	} else {
		if phoneNumber == "" {
			return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INVALID_ERROR, "both id and phonenumber cannot m=be missing")
		}
		params.AccountNumber = sql.NullString{
			Valid: true,
			String: phoneNumber,
		}
		params2.AccountNumber = sql.NullString{
			Valid: true,
			String: phoneNumber,
		}
	}

	nonPosteds, err := r.queries.GetClientsNonPosted(tc, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get non-posted: %s", err.Error())
		} 

		setSpanError(span, codes.Error, err, "failed to get client data")
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get non-posted: %s", err.Error())	
	}

	totalNonPosted, err := r.queries.CountClientsNonPosted(ctx, params2)
	if err != nil {
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get total loans: %s", err)
	}

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

	return paymentDetails, pkg.CreatePaginationMetadata(uint32(totalNonPosted), pgData.PageSize, pgData.CurrentPage), nil
}

// func (r *HelperRepository)userToClientDashboard(ctx context.Context, id uint32) repository.ClientDashboardResponse {
// 	client, _ := r.queries.GetClient(ctx, id)

// 	branch, _ := r.queries.GetBranch(ctx, client.BranchID)

// 	return repository.ClientDashboardResponse{
// 		ID: client.ID,
// 		FullName: client.FullName,
// 		PhoneNumber: client.PhoneNumber,
// 		IdNumber: client.IDNumber.String,
// 		Dob: client.Dob.Time.String(),
// 		Gender: string(client.Gender),
// 		Active: client.Active,
// 		Overpayment: client.Overpayment,
// 		CreatedAt: client.CreatedAt,
// 		AssignedStaff: r.clientToUserDashboard(ctx, client.AssignedStaff),
// 		CreatedBy: r.clientToUserDashboard(ctx, client.CreatedBy),
// 		BranchName: branch.Name,
// 	}
// } 

// func (r *HelperRepository)clientToUserDashboard(ctx context.Context, id uint32) repository.UserDashboardResponse {
// 	user, _ := r.queries.GetUser(ctx, id)

// 	branch, _ := r.queries.GetBranch(ctx, user.BranchID)

// 	return repository.UserDashboardResponse{
// 		ID: user.ID,
// 		Fullname: user.FullName,
// 		Email: user.Email,
// 		PhoneNumber: user.PhoneNumber,
// 		Role: string(user.Role),
// 		BranchName: branch.Name,
// 		CreatedAt: user.CreatedAt,
// 	}
// }