package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
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

func (r *LoanRepository) TransferLoan(ctx context.Context, officerId uint32, loanId uint32, adminId uint32) error {
	_, err := r.queries.TransferLoan(ctx, generated.TransferLoanParams{
		ID:          loanId,
		LoanOfficer: officerId,
		UpdatedBy: sql.NullInt32{
			Valid: true,
			Int32: int32(adminId),
		},
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

func (t *LoanRepository) GetClientActiceLoan(ctx context.Context, clientID uint32) (uint32, error) {
	loanID, err := t.queries.GetClientActiveLoan(ctx, generated.GetClientActiveLoanParams{
		ClientID: clientID,
		Status:   generated.LoansStatusACTIVE,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, pkg.Errorf(pkg.NOT_FOUND_ERROR, "loan not found")
		}

		return 0, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loan: %s", err.Error())
	}

	return loanID, nil
}

func (r *LoanRepository) ListLoans(ctx context.Context, category *repository.Category, pgData *pkg.PaginationMetadata) ([]repository.Loan, pkg.PaginationMetadata, error) {
	params := generated.ListLoansParams{
		Column1: "",
		FullName:   "", 
		FullName_2: "", 
		Column4: "",
		FINDINSET: "", 
		Limit:    int32(pgData.PageSize),
		Offset:   int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
	}

	params2 := generated.CountLoansParams{
		Column1: "",
		FullName:   "", 
		FullName_2: "", 
		Column4: "",
		FINDINSET: "", 
	}

	if category.Search != nil {
		searchValue := "%" + *category.Search + "%"
		params.Column1 = "has_search"
		params.FullName = searchValue
		params.FullName_2 = searchValue

		params2.Column1 = "has_search"
		params2.FullName = searchValue
		params2.FullName_2 = searchValue
	}

	if category.Statuses != nil {
		params.Column4 = "has_status"
		params2.Column4 = "has_status"		
		params.FINDINSET = *category.Statuses
		params2.FINDINSET = *category.Statuses
	}

	loans, err := r.queries.ListLoans(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loans found")
		}
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loans: %s", err.Error())
	}

	totalLoans, err := r.queries.CountLoans(ctx, params2)
	if err != nil {
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get total loans: %s", err.Error())
	}

	result := make([]repository.Loan, len(loans))

	for i, loan := range loans {
		result[i] = convertGeneratedLoan(generated.Loan{
			ID:                 loan.ID,
			ProductID:          loan.ProductID,
			ClientID:           loan.ClientID,
			LoanOfficer:        loan.LoanOfficer,
			LoanPurpose:        loan.LoanPurpose,
			DueDate:            loan.DueDate,
			ApprovedBy:         loan.ApprovedBy,
			DisbursedOn:        loan.DisbursedOn,
			DisbursedBy:        loan.DisbursedBy,
			TotalInstallments:  loan.TotalInstallments,
			InstallmentsPeriod: loan.InstallmentsPeriod,
			Status:             loan.Status,
			ProcessingFee:      loan.ProcessingFee,
			PaidAmount:         loan.PaidAmount,
			UpdatedBy:          loan.UpdatedBy,
			CreatedBy:          loan.CreatedBy,
			CreatedAt:          loan.CreatedAt,
			FeePaid:            loan.FeePaid,
		})
	}

	return result, pkg.CreatePaginationMetadata(uint32(totalLoans), pgData.PageSize, pgData.CurrentPage), nil
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

func (r *LoanRepository) GetLoanInstallments(ctx context.Context, id uint32, pgData *pkg.PaginationMetadata) ([]repository.Installment, error) {
	installments, err := r.queries.ListInstallmentsByLoan(ctx, generated.ListInstallmentsByLoanParams{
		LoanID: id,
		Limit: int32(pgData.PageSize),
		Offset: int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
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
		FeePaid:            loan.FeePaid,
	}
}

func (r *LoanRepository) GetReportLoanData(ctx context.Context, filters services.ReportFilters) ([]services.LoanReportData, error) {
	loans, err := r.GetLoansReportData(ctx, GetLoansReportDataParams{
		StartDate: filters.StartDate,
		EndDate: filters.EndDate,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loans found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get report loan data: %s", err.Error())
	}

	rslt := make([]services.LoanReportData, len(loans))

	for i, loan := range loans {
		rslt[i] = services.LoanReportData{
			LoanID: loan.LoanID,
			ClientName: loan.ClientName,
			BranchName: loan.BranchName,
			LoanOfficer: loan.LoanOfficer,
			LoanAmount: loan.LoanAmount,
			RepayAmount: loan.RepayAmount,
			PaidAmount: loan.PaidAmount,
			OutstandingAmount: pkg.InterfaceFloat64(loan.OutstandingAmount),
			Status: loan.Status,
			TotalInstallments: loan.TotalInstallments,
			PaidInstallments: loan.PaidInstallments,
			DueDate: pkg.TimePtr(loan.DueDate.Time),
			DisbursedDate: pkg.TimePtr(loan.DisbursedDate.Time),
			DefaultRisk: pkg.InterfaceFloat64(loan.DefaultRisk),
		}
	}

	return rslt, nil
}

func (r *LoanRepository) GetReportLoanByIdData(ctx context.Context,id uint32) (services.LoanReportDataById, error) {
	loan, err := r.GetLoanReportDataById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return services.LoanReportDataById{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loans found")
		}

		return services.LoanReportDataById{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get report loan by id data: %s", err.Error())
	}

	return convertLoanReportDataById(loan)
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

func convertLoanReportDataById(row GetLoanReportDataByIdRow) (services.LoanReportDataById, error) {
	var installments []services.LoanReportDataByIdInstallmentDetails
	if row.InstallmentDetails != nil {
		installmentsBytes, ok := row.InstallmentDetails.([]byte)
		if !ok {
			return services.LoanReportDataById{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to convert installments to bytes")
		}

		err := json.Unmarshal(installmentsBytes, &installments) 
		if err != nil {
			return services.LoanReportDataById{}, pkg.Errorf(pkg.INTERNAL_ERROR, "error unmarshalling installments: ", err)
		}
	}

	return services.LoanReportDataById{
		LoanID:                row.LoanID,
		ClientName:            row.ClientName,
		LoanAmount:            row.LoanAmount,
		RepayAmount:           row.RepayAmount,
		PaidAmount:            row.PaidAmount,
		Status:                row.Status,
		TotalInstallments:     row.TotalInstallments,
		PaidInstallments:      row.PaidInstallments,
		RemainingInstallments: row.RemainingInstallments,
		InstallmentDetails:    installments,
	}, nil
}
