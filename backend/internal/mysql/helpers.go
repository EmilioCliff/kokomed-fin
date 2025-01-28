package mysql

import (
	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
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

func (r *HelperRepository) GetDashboardData() (repository.DashboardData, error) {
	// sql command to get inactive loans
	// sql command to get recent payments
	// sql command to get widget data:
		// no of total customers, active customers and inactive customers
		// no of total loans, active loans and inactive loans
		// no of total transactions, no of disbursed loans and number of repaid loans
		// no of all payments, no of posted_payments(transactions linked to client) and non_posted payments
	panic("not implemented") // TODO: Implement
}

func (r *HelperRepository) GetProductData() ([]repository.ProductData, error) {
	// sql command to get all loans
	panic("not implemented") // TODO: Implement
}
func (r *HelperRepository) GetClientData() ([]repository.ClientData, error) {
	// sql command to get all clients
	panic("not implemented") // TODO: Implement
}
func (r *HelperRepository) GetLoanOfficerData() ([]repository.LoanOfficerData, error) {
	// sql command to get all loan officers
	panic("not implemented") // TODO: Implement
}