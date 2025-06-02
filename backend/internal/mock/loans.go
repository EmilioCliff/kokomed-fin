package mock

import (
	"context"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.LoansRepository = (*MockLoanRepository)(nil)

type MockLoanRepository struct {
	MockCreateLoan                 func(ctx context.Context, loan *repository.Loan) (repository.LoanFullData, error)
	MockDisburseLoan               func(ctx context.Context, disburseLoan *repository.DisburseLoan) (uint32, error)
	MockTransferLoan               func(ctx context.Context, officerId uint32, loanId uint32, adminId uint32) error
	MockGetLoanByID                func(ctx context.Context, id uint32) (repository.Loan, error)
	MockGetClientActiceLoan        func(ctx context.Context, clientID uint32) (uint32, error)
	MockDeleteLoan                 func(ctx context.Context, id uint32) error
	MockGetLoanInstallments        func(ctx context.Context, id uint32) ([]repository.Installment, error)
	MockGetInstallment             func(ctx context.Context, id uint32) (repository.Installment, error)
	MockUpdateInstallment          func(ctx context.Context, installment *repository.UpdateInstallment) (repository.Installment, error)
	MockGetReportLoanByIdData      func(ctx context.Context, id uint32) (services.LoanReportDataById, error)
	MockGetLoanEvents              func(ctx context.Context) ([]repository.LoanEvent, error)
	MockListLoans                  func(ctx context.Context, category *repository.Category, pgData *pkg.PaginationMetadata) ([]repository.LoanFullData, pkg.PaginationMetadata, error)
	MockGetExpectedPayments        func(ctx context.Context, category *repository.Category, pgData *pkg.PaginationMetadata) ([]repository.ExpectedPayment, pkg.PaginationMetadata, error)
	MockListUnpaidInstallmentsData func(ctx context.Context, category *repository.Category, pgData *pkg.PaginationMetadata) ([]repository.UnpaidInstallmentData, pkg.PaginationMetadata, error)
	MockGetReportLoanData          func(ctx context.Context, filters services.ReportFilters) ([]services.LoanReportData, services.LoanSummary, error)
}

func (m *MockLoanRepository) CreateLoan(
	ctx context.Context,
	loan *repository.Loan,
) (repository.LoanFullData, error) {
	return m.MockCreateLoan(ctx, loan)
}

func (m *MockLoanRepository) DisburseLoan(
	ctx context.Context,
	disburseLoan *repository.DisburseLoan,
) (uint32, error) {
	return m.MockDisburseLoan(ctx, disburseLoan)
}

func (m *MockLoanRepository) TransferLoan(
	ctx context.Context,
	officerId uint32,
	loanId uint32,
	adminId uint32,
) error {
	return m.MockTransferLoan(ctx, officerId, loanId, adminId)
}
func (m *MockLoanRepository) GetLoanByID(ctx context.Context, id uint32) (repository.Loan, error) {
	return m.MockGetLoanByID(ctx, id)
}

func (m *MockLoanRepository) GetClientActiceLoan(
	ctx context.Context,
	clientID uint32,
) (uint32, error) {
	return m.MockGetClientActiceLoan(ctx, clientID)
}
func (m *MockLoanRepository) DeleteLoan(ctx context.Context, id uint32) error {
	return m.MockDeleteLoan(ctx, id)
}

func (m *MockLoanRepository) GetLoanInstallments(
	ctx context.Context,
	id uint32,
) ([]repository.Installment, error) {
	return m.MockGetLoanInstallments(ctx, id)
}

func (m *MockLoanRepository) GetInstallment(
	ctx context.Context,
	id uint32,
) (repository.Installment, error) {
	return m.MockGetInstallment(ctx, id)
}

func (m *MockLoanRepository) UpdateInstallment(
	ctx context.Context,
	installment *repository.UpdateInstallment,
) (repository.Installment, error) {
	return m.MockUpdateInstallment(ctx, installment)
}

func (m *MockLoanRepository) GetReportLoanByIdData(
	ctx context.Context,
	id uint32,
) (services.LoanReportDataById, error) {
	return m.MockGetReportLoanByIdData(ctx, id)
}
func (m *MockLoanRepository) GetLoanEvents(ctx context.Context) ([]repository.LoanEvent, error) {
	return m.MockGetLoanEvents(ctx)
}

func (m *MockLoanRepository) ListLoans(
	ctx context.Context,
	category *repository.Category,
	pgData *pkg.PaginationMetadata,
) ([]repository.LoanFullData, pkg.PaginationMetadata, error) {
	return m.MockListLoans(ctx, category, pgData)
}

func (m *MockLoanRepository) GetExpectedPayments(
	ctx context.Context,
	category *repository.Category,
	pgData *pkg.PaginationMetadata,
) ([]repository.ExpectedPayment, pkg.PaginationMetadata, error) {
	return m.MockGetExpectedPayments(ctx, category, pgData)
}

func (m *MockLoanRepository) ListUnpaidInstallmentsData(
	ctx context.Context,
	category *repository.Category,
	pgData *pkg.PaginationMetadata,
) ([]repository.UnpaidInstallmentData, pkg.PaginationMetadata, error) {
	return m.MockListUnpaidInstallmentsData(ctx, category, pgData)
}

func (m *MockLoanRepository) GetReportLoanData(
	ctx context.Context,
	filters services.ReportFilters,
) ([]services.LoanReportData, services.LoanSummary, error) {
	return m.MockGetReportLoanData(ctx, filters)
}
