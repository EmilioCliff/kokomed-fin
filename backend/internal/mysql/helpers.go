package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
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
	widgetsData, err := r.queries.DashBoardDataHelper(ctx)
	if err != nil {
		return repository.DashboardData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting dashboard widgets")
	}

	loansData, err := r.queries.DashBoardInactiveLoans(ctx)
	if err != nil {
		return repository.DashboardData{}, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting dashboard inactive loans")
	}

	paymentsData, err := r.queries.DashBoardRecentsPayments(ctx)
	if err != nil {
		return repository.DashboardData{}, pkg.Errorf(pkg.INTERNAL_ERROR, err.Error())
	}

	inactiveLoans := make([]repository.InactiveLoan, len(loansData))
	for idx, loans := range loansData {
		inactiveLoans[idx] = repository.InactiveLoan{
			ID:          loans.ID,
			Amount:      loans.LoanAmount,
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

		totalLoanAmountBtye, _ := widgetsData.TotalLoanAmount.([]byte)
		totalLoanAmount, _ := strconv.ParseFloat(string(totalLoanAmountBtye), 64)
		totalLoanDisbursedBtye, _ := widgetsData.TotalLoanDisbursed.([]byte)
		totalLoanDisbursed, _ := strconv.ParseFloat(string(totalLoanDisbursedBtye), 64)
		totalLoanPaidBtye, _ := widgetsData.TotalLoanPaid.([]byte)
		totalLoanPaid, _ := strconv.ParseFloat(string(totalLoanPaidBtye), 64)
		totalPaymentsReceivedBtye, _ := widgetsData.TotalPaymentsReceived.([]byte)
		totalPaymentsReceived, _ := strconv.ParseFloat(string(totalPaymentsReceivedBtye), 64)
		totalNonPostedBtye, _ := widgetsData.TotalNonPosted.([]byte)
		totalNonPosted, _ := strconv.ParseFloat(string(totalNonPostedBtye), 64)

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
				MainAmount: totalLoanAmount,
				Active: totalLoanDisbursed,
				ActiveTitle: "Disbursed",
				Closed: totalLoanPaid,
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
	products, err := r.queries.HelperProduct(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting product list data")		
	}

	rsp := make([]repository.ProductData, len(products))
	for _, product := range products {
		rsp = append(rsp, repository.ProductData{
			ID: product.Productid,
			Name: fmt.Sprintf("%f %s", product.Loanamount, product.Branchname),
		})
	}

	return rsp, nil
}
func (r *HelperRepository) GetClientData(ctx context.Context) ([]repository.ClientData, error) {
	clients, err := r.queries.HelperClient(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting clients list data")		
	}

	rsp := make([]repository.ClientData, len(clients))
	for _, client := range clients {
		rsp = append(rsp, repository.ClientData{
			ID: client.ID,
			Name: client.FullName,
		})
	}

	return rsp, nil
}
func (r *HelperRepository) GetLoanOfficerData(ctx context.Context) ([]repository.LoanOfficerData, error) {
	users, err := r.queries.HelperUser(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting user list data")		
	}

	rsp := make([]repository.LoanOfficerData, len(users))
	for _, user := range users {
		rsp = append(rsp, repository.LoanOfficerData{
			ID: user.ID,
			Name: user.FullName,
		})
	}

	return rsp, nil
}

func (r *HelperRepository)userToClientDashboard(ctx context.Context, id uint32) repository.ClientDashboardResponse {
	client, _ := r.queries.GetClient(ctx, id)

	branch, _ := r.queries.GetBranch(ctx, client.BranchID)

	return repository.ClientDashboardResponse{
		ID: client.ID,
		FullName: client.FullName,
		PhoneNumber: client.PhoneNumber,
		IdNumber: client.IDNumber.String,
		Dob: client.Dob.Time.String(),
		Gender: string(client.Gender),
		Active: client.Active,
		Overpayment: client.Overpayment,
		CreatedAt: client.CreatedAt,
		AssignedStaff: r.clientToUserDashboard(ctx, client.AssignedStaff),
		CreatedBy: r.clientToUserDashboard(ctx, client.CreatedBy),
		BranchName: branch.Name,
	}
} 

func (r *HelperRepository)clientToUserDashboard(ctx context.Context, id uint32) repository.UserDashboardResponse {
	user, _ := r.queries.GetUser(ctx, id)

	branch, _ := r.queries.GetBranch(ctx, user.BranchID)

	return repository.UserDashboardResponse{
		ID: user.ID,
		Fullname: user.FullName,
		Email: user.Email,
		PhoneNumber: user.PhoneNumber,
		Role: string(user.Role),
		BranchName: branch.Name,
		CreatedAt: user.CreatedAt,
	}
}