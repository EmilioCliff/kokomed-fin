package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.LoansRepository = (*LoanRepository)(nil)

type LoanRepository struct {
	db      *Store
	queries generated.Querier
}

func NewLoanRepository(db *Store) *LoanRepository {
	return &LoanRepository{
		db:      db,
		queries: generated.New(db.db),
	}
}

func (r *LoanRepository) CreateLoan(ctx context.Context, loan *repository.Loan) (repository.Loan, error) {
	params := generated.CreateLoanParams{
		ProductID:          loan.ProductID,
		ClientID:           loan.ClientID,
		LoanOfficer:        loan.LoanOfficerID,
		ApprovedBy:         loan.ApprovedBy,
		TotalInstallments:  loan.TotalInstallments,
		InstallmentsPeriod: loan.InstallmentsPeriod,
		ProcessingFee:      loan.ProcessingFee,
		CreatedBy:          loan.CreatedBy,
	}

	if loan.LoanPurpose != nil {
		params.LoanPurpose = sql.NullString{
			Valid:  true,
			String: *loan.LoanPurpose,
		}
	}

	if loan.DueDate != nil && loan.DisbursedBy != nil && loan.DisbursedOn != nil {
		params.DueDate = sql.NullTime{
			Valid: true,
			Time:  *loan.DueDate,
		}
		params.DisbursedOn = sql.NullTime{
			Valid: true,
			Time:  *loan.DisbursedOn,
		}
		params.DisbursedBy = sql.NullInt32{
			Valid: true,
			Int32: int32(*loan.DisbursedBy),
		}
	}

	execResult, err := r.queries.CreateLoan(ctx, params)
	if err != nil {
		return repository.Loan{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create loan: %s", err.Error())
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.Loan{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	loan.ID = uint32(id)

	return *loan, nil
}

func (r *LoanRepository) DisburseLoan(ctx context.Context, disburseLoan *repository.DisburseLoan) error {
	_, err := r.queries.DisburseLoan(ctx, generated.DisburseLoanParams{
		ID: disburseLoan.ID,
		DisbursedOn: sql.NullTime{
			Valid: true,
			Time:  disburseLoan.DisbursedOn,
		},
		DisbursedBy: sql.NullInt32{
			Valid: true,
			Int32: int32(disburseLoan.DisbursedBy),
		},
		DueDate: sql.NullTime{
			Valid: true,
			Time:  disburseLoan.DueDate,
		},
	})
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to disburse loan: %s", err.Error())
	}

	return nil
}

func (r *LoanRepository) UpdateLoan(ctx context.Context, loan *repository.UpdateLoan) (repository.Loan, error) {
	params := generated.UpdateLoanParams{
		ID:         loan.ID,
		PaidAmount: loan.PaidAmount,
	}

	if loan.UpdatedBy != nil {
		params.UpdatedBy = sql.NullInt32{
			Valid: true,
			Int32: int32(*loan.UpdatedBy),
		}
	}

	execResult, err := r.queries.UpdateLoan(ctx, params)
	if err != nil {
		return repository.Loan{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update loan: %s", err.Error())
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.Loan{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	return r.GetLoanByID(ctx, uint32(id))
}

func (r *LoanRepository) TransferLoan(ctx context.Context, officerId uint32, loanId uint32) error {
	_, err := r.queries.TransferLoan(ctx, generated.TransferLoanParams{
		ID:          loanId,
		LoanOfficer: officerId,
	})
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to transfer loan: %s", err.Error())
	}

	return nil
}

func (r *LoanRepository) GetLoanByID(ctx context.Context, id uint32) (repository.Loan, error) {
	loan, err := r.queries.GetLoan(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.Loan{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "loan not found")
		}

		return repository.Loan{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loan: %s", err.Error())
	}

	return convertGeneratedLoan(loan), nil
}

func (r *LoanRepository) ListLoans(ctx context.Context, pgData *pkg.PaginationMetadata) ([]repository.Loan, error) {
	loans, err := r.queries.ListLoans(ctx, generated.ListLoansParams{
		Limit:  pkg.GetPageSize(),
		Offset: pkg.CalculateOffset(pgData.CurrentPage),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loans found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loans: %s", err.Error())
	}

	result := make([]repository.Loan, len(loans))

	for i, loan := range loans {
		result[i] = convertGeneratedLoan(loan)
	}

	return result, nil
}

func (r *LoanRepository) DeleteLoan(ctx context.Context, id uint32) error {
	err := r.queries.DeleteLoan(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loan found")
		}

		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to delete loan: %s", err.Error())
	}

	return nil
}

func (r *LoanRepository) GetAllClientsLoans(ctx context.Context, id uint32, pgData *pkg.PaginationMetadata) ([]repository.Loan, error) {
	loans, err := r.queries.ListLoansByClient(ctx, generated.ListLoansByClientParams{
		ClientID: id,
		Limit:    pkg.GetPageSize(),
		Offset:   pkg.CalculateOffset(pgData.CurrentPage),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loans found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loans: %s", err.Error())
	}

	result := make([]repository.Loan, len(loans))

	for i, loan := range loans {
		result[i] = convertGeneratedLoan(loan)
	}

	return result, nil
}

func (r *LoanRepository) GetAllUsersLoans(ctx context.Context, id uint32, pgData *pkg.PaginationMetadata) ([]repository.Loan, error) {
	loans, err := r.queries.ListLoansByLoanOfficer(ctx, generated.ListLoansByLoanOfficerParams{
		LoanOfficer: id,
		Limit:       pkg.GetPageSize(),
		Offset:      pkg.CalculateOffset(pgData.CurrentPage),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loans found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loans: %s", err.Error())
	}

	result := make([]repository.Loan, len(loans))

	for i, loan := range loans {
		result[i] = convertGeneratedLoan(loan)
	}

	return result, nil
}

func (r *LoanRepository) ListLoansByStatus(ctx context.Context, status string, pgData *pkg.PaginationMetadata) ([]repository.Loan, error) {
	loans, err := r.queries.ListLoansByStatus(ctx, generated.ListLoansByStatusParams{
		Status: generated.LoansStatus(status),
		Limit:  pkg.GetPageSize(),
		Offset: pkg.CalculateOffset(pgData.CurrentPage),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loans found with status: %s", status)
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loans: %s", err.Error())
	}

	result := make([]repository.Loan, len(loans))

	for i, loan := range loans {
		result[i] = convertGeneratedLoan(loan)
	}

	return result, nil
}

func (r *LoanRepository) ListNonDisbursedLoans(ctx context.Context, pgData *pkg.PaginationMetadata) ([]repository.Loan, error) {
	loans, err := r.queries.ListNonDisbursedLoans(ctx, generated.ListNonDisbursedLoansParams{
		Limit:  pkg.GetPageSize(),
		Offset: pkg.CalculateOffset(pgData.CurrentPage),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loans found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loans: %s", err.Error())
	}

	result := make([]repository.Loan, len(loans))

	for i, loan := range loans {
		result[i] = convertGeneratedLoan(loan)
	}

	return result, nil
}

func (r *LoanRepository) UpdateLoanStatus(ctx context.Context, id uint32, status string) error {
	_, err := r.queries.UpdateLoanStatus(ctx, generated.UpdateLoanStatusParams{
		Status: generated.LoansStatus(status),
		ID:     id,
	})
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to change loans status: %s", err.Error())
	}

	return nil
}

func (r *LoanRepository) CreateInstallment(ctx context.Context, installment *repository.Installment) (repository.Installment, error) {
	execResult, err := r.queries.CreateInstallment(ctx, generated.CreateInstallmentParams{
		LoanID:            installment.LoanID,
		InstallmentNumber: installment.InstallmentNo,
		AmountDue:         installment.Amount,
		RemainingAmount:   installment.RemainingAmount,
		DueDate:           installment.DueDate,
	})
	if err != nil {
		return repository.Installment{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create installment: %v", err)
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.Installment{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	installment.ID = uint32(id)

	return *installment, err
}

func (r *LoanRepository) GetLoanInstallments(ctx context.Context, id uint32, pgData *pkg.PaginationMetadata) ([]repository.Installment, error) {
	installments, err := r.queries.ListInstallmentsByLoan(ctx, generated.ListInstallmentsByLoanParams{
		LoanID: id,
		Limit:  pkg.GetPageSize(),
		Offset: pkg.CalculateOffset(pgData.CurrentPage),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loans installments found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loans installments: %s", err.Error())
	}

	result := make([]repository.Installment, len(installments))

	for i, installment := range installments {
		result[i] = convertGeneratedInstallment(installment)
	}

	return result, nil
}

func (r *LoanRepository) GetInstallment(ctx context.Context, id uint32) (repository.Installment, error) {
	installment, err := r.queries.GetInstallment(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.Installment{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no installment found")
		}

		return repository.Installment{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get intallment: %s", err.Error())
	}

	return convertGeneratedInstallment(installment), nil
}

func (r *LoanRepository) UpdateInstallment(ctx context.Context, installment *repository.UpdateInstallment) (repository.Installment, error) {
	params := generated.UpdateInstallmentParams{
		ID:              installment.ID,
		RemainingAmount: installment.RemainingAmount,
	}

	if installment.Paid != nil {
		params.Paid = sql.NullBool{
			Valid: true,
			Bool:  *installment.Paid,
		}
	}

	if installment.PaidAt != nil {
		params.PaidAt = sql.NullTime{
			Valid: true,
			Time:  *installment.PaidAt,
		}
	}

	execResult, err := r.queries.UpdateInstallment(ctx, params)
	if err != nil {
		return repository.Installment{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update installment: %s", err.Error())
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.Installment{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	return r.GetInstallment(ctx, uint32(id))
}

func convertGeneratedLoan(loan generated.Loan) repository.Loan {
	var loanPurpose *string
	if loan.LoanPurpose.Valid {
		loanPurpose = &loan.LoanPurpose.String
	}

	var dueDate *time.Time
	if loan.DueDate.Valid {
		dueDate = &loan.DueDate.Time
	}

	var disbursedOn *time.Time
	if loan.DisbursedOn.Valid {
		disbursedOn = &loan.DisbursedOn.Time
	}

	var disbursedBy *uint32

	if loan.DisbursedBy.Valid {
		value := uint32(loan.DisbursedBy.Int32)
		disbursedBy = &value
	}

	var updatedBy *uint32

	if loan.UpdatedBy.Valid {
		value := uint32(loan.UpdatedBy.Int32)
		updatedBy = &value
	}

	return repository.Loan{
		ID:                 loan.ID,
		ProductID:          loan.ProductID,
		ClientID:           loan.ClientID,
		LoanOfficerID:      loan.LoanOfficer,
		LoanPurpose:        loanPurpose,
		DueDate:            dueDate,
		ApprovedBy:         loan.ApprovedBy,
		DisbursedOn:        disbursedOn,
		DisbursedBy:        disbursedBy,
		TotalInstallments:  loan.TotalInstallments,
		InstallmentsPeriod: loan.InstallmentsPeriod,
		Status:             string(loan.Status),
		ProcessingFee:      loan.ProcessingFee,
		PaidAmount:         loan.PaidAmount,
		UpdatedBy:          updatedBy,
		CreatedBy:          loan.CreatedBy,
		CreatedAt:          loan.CreatedAt,
	}
}

func convertGeneratedInstallment(installment generated.Installment) repository.Installment {
	return repository.Installment{
		ID:              installment.ID,
		LoanID:          installment.LoanID,
		InstallmentNo:   installment.InstallmentNumber,
		Amount:          installment.AmountDue,
		RemainingAmount: installment.RemainingAmount,
		Paid:            installment.Paid,
		PaidAt:          installment.PaidAt.Time,
		DueDate:         installment.DueDate,
	}
}
