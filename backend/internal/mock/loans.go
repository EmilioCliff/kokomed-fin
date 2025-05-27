package mock

import (
	"context"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.LoansRepository = (*MockLoanRepository)(nil)

type MockLoanRepository struct {
	mockCreateLoan                 func(ctx context.Context, loan *repository.Loan) (repository.LoanFullData, error)
	mockDisburseLoan               func(ctx context.Context, disburseLoan *repository.DisburseLoan) (uint32, error)
	mockTransferLoan               func(ctx context.Context, officerId uint32, loanId uint32, adminId uint32) error
	mockGetLoanByID                func(ctx context.Context, id uint32) (repository.Loan, error)
	mockGetClientActiceLoan        func(ctx context.Context, clientID uint32) (uint32, error)
	mockDeleteLoan                 func(ctx context.Context, id uint32) error
	mockGetLoanInstallments        func(ctx context.Context, id uint32) ([]repository.Installment, error)
	mockGetInstallment             func(ctx context.Context, id uint32) (repository.Installment, error)
	mockUpdateInstallment          func(ctx context.Context, installment *repository.UpdateInstallment) (repository.Installment, error)
	mockGetReportLoanByIdData      func(ctx context.Context, id uint32) (services.LoanReportDataById, error)
	mockGetLoanEvents              func(ctx context.Context) ([]repository.LoanEvent, error)
	mockListLoans                  func(ctx context.Context, category *repository.Category, pgData *pkg.PaginationMetadata) ([]repository.LoanFullData, pkg.PaginationMetadata, error)
	mockGetExpectedPayments        func(ctx context.Context, category *repository.Category, pgData *pkg.PaginationMetadata) ([]repository.ExpectedPayment, pkg.PaginationMetadata, error)
	mockListUnpaidInstallmentsData func(ctx context.Context, category *repository.Category, pgData *pkg.PaginationMetadata) ([]repository.UnpaidInstallmentData, pkg.PaginationMetadata, error)
	mockGetReportLoanData          func(ctx context.Context, filters services.ReportFilters) ([]services.LoanReportData, services.LoanSummary, error)
}

func (m *MockLoanRepository) CreateLoan(
	ctx context.Context,
	loan *repository.Loan,
) (repository.LoanFullData, error) {
	return m.mockCreateLoan(ctx, loan)
}

func (m *MockLoanRepository) DisburseLoan(
	ctx context.Context,
	disburseLoan *repository.DisburseLoan,
) (uint32, error) {
	return m.mockDisburseLoan(ctx, disburseLoan)
}

func (m *MockLoanRepository) TransferLoan(
	ctx context.Context,
	officerId uint32,
	loanId uint32,
	adminId uint32,
) error {
	return m.mockTransferLoan(ctx, officerId, loanId, adminId)
}
func (m *MockLoanRepository) GetLoanByID(ctx context.Context, id uint32) (repository.Loan, error) {
	return m.mockGetLoanByID(ctx, id)
}

func (m *MockLoanRepository) GetClientActiceLoan(
	ctx context.Context,
	clientID uint32,
) (uint32, error) {
	return m.mockGetClientActiceLoan(ctx, clientID)
}
func (m *MockLoanRepository) DeleteLoan(ctx context.Context, id uint32) error {
	return m.mockDeleteLoan(ctx, id)
}

func (m *MockLoanRepository) GetLoanInstallments(
	ctx context.Context,
	id uint32,
) ([]repository.Installment, error) {
	return m.mockGetLoanInstallments(ctx, id)
}

func (m *MockLoanRepository) GetInstallment(
	ctx context.Context,
	id uint32,
) (repository.Installment, error) {
	return m.mockGetInstallment(ctx, id)
}

func (m *MockLoanRepository) UpdateInstallment(
	ctx context.Context,
	installment *repository.UpdateInstallment,
) (repository.Installment, error) {
	return m.mockUpdateInstallment(ctx, installment)
}

func (m *MockLoanRepository) GetReportLoanByIdData(
	ctx context.Context,
	id uint32,
) (services.LoanReportDataById, error) {
	return m.mockGetReportLoanByIdData(ctx, id)
}
func (m *MockLoanRepository) GetLoanEvents(ctx context.Context) ([]repository.LoanEvent, error) {
	return m.mockGetLoanEvents(ctx)
}

func (m *MockLoanRepository) ListLoans(
	ctx context.Context,
	category *repository.Category,
	pgData *pkg.PaginationMetadata,
) ([]repository.LoanFullData, pkg.PaginationMetadata, error) {
	return m.mockListLoans(ctx, category, pgData)
}

func (m *MockLoanRepository) GetExpectedPayments(
	ctx context.Context,
	category *repository.Category,
	pgData *pkg.PaginationMetadata,
) ([]repository.ExpectedPayment, pkg.PaginationMetadata, error) {
	return m.mockGetExpectedPayments(ctx, category, pgData)
}

func (m *MockLoanRepository) ListUnpaidInstallmentsData(
	ctx context.Context,
	category *repository.Category,
	pgData *pkg.PaginationMetadata,
) ([]repository.UnpaidInstallmentData, pkg.PaginationMetadata, error) {
	return m.mockListUnpaidInstallmentsData(ctx, category, pgData)
}

func (m *MockLoanRepository) GetReportLoanData(
	ctx context.Context,
	filters services.ReportFilters,
) ([]services.LoanReportData, services.LoanSummary, error) {
	return m.mockGetReportLoanData(ctx, filters)
}
