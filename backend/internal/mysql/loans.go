package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

func (r *LoanRepository) GetExpectedPayments(ctx context.Context, category *repository.Category, pgData *pkg.PaginationMetadata) ([]repository.ExpectedPayment, pkg.PaginationMetadata, error) {
	params := generated.ListExpectedPaymentsParams{
		Column1: "",
		FullName:   "", 
		FullName_2: "", 
		Limit:    int32(pgData.PageSize),
		Offset:   int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
	}

	params2 := generated.CountExpectedPaymentsParams{
		Column1: "",
		FullName:   "", 
		FullName_2: "", 
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

	expectedPayments, err := r.queries.ListExpectedPayments(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no unexpected payments found")
		}
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get unexpected payments: %s", err.Error())
	}

	totalPayments, err := r.queries.CountExpectedPayments(ctx, params2)
	if err != nil {
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get total unexpected payments: %s", err.Error())
	}

	result := make([]repository.ExpectedPayment, len(expectedPayments))

	for i, payment := range expectedPayments {
		result[i] = repository.ExpectedPayment{
			LoanId: payment.LoanID,
			BranchName: payment.BranchName,
			ClientName: payment.ClientName,
			LoanOfficerName: payment.LoanOfficerName,
			LoanAmount: payment.LoanAmount,
			RepayAmount: payment.RepayAmount,
			TotalUnpaid: pkg.InterfaceFloat64(payment.TotalUnpaid),
			DueDate: payment.DueDate.Time.Format("2006-02-01"),
		}
	}

	return result, pkg.CreatePaginationMetadata(uint32(totalPayments), pgData.PageSize, pgData.CurrentPage), nil
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

func (r *LoanRepository) GetLoanInstallments(ctx context.Context, id uint32) ([]repository.Installment, error) {
	installments, err := r.queries.ListInstallmentsByLoan(ctx, id)
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


func (r *LoanRepository) GetReportLoanData(ctx context.Context, filters services.ReportFilters) ([]services.LoanReportData, services.LoanSummary, error) {
	loans, err := r.GetLoansReportData(ctx, GetLoansReportDataParams{
		StartDate: filters.StartDate,
		EndDate:   filters.EndDate,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, services.LoanSummary{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loans found")
		}
		return nil, services.LoanSummary{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get report loan data: %s", err.Error())
	}

	rslt := make([]services.LoanReportData, len(loans))

	summary := services.LoanSummary{}
	branchLoanCount := map[string]int64{}
	officerLoanCount := map[string]int64{}

	for i, loan := range loans {
		dueDate := ""
		disbursedDate := ""
		if loan.DueDate.Valid {
			dueDate = loan.DueDate.Time.Format("2006-01-02")
		}
		if loan.DisbursedDate.Valid {
			disbursedDate = loan.DisbursedDate.Time.Format("2006-01-02")
		}

		rslt[i] = services.LoanReportData{
			LoanID:            loan.LoanID,
			ClientName:        loan.ClientName,
			BranchName:        loan.BranchName,
			LoanOfficer:       loan.LoanOfficer,
			LoanAmount:        loan.LoanAmount,
			RepayAmount:       loan.RepayAmount,
			PaidAmount:        loan.PaidAmount,
			OutstandingAmount: loan.RepayAmount - loan.PaidAmount,
			Status:            loan.Status,
			TotalInstallments: loan.TotalInstallments,
			PaidInstallments:  loan.PaidInstallments,
			DueDate:           dueDate,
			DisbursedDate:     disbursedDate,
			DefaultRisk: pkg.InterfaceFloat64(loan.DefaultRisk),
		}

		summary.TotalLoans++
		summary.TotalDisbursedAmount += loan.LoanAmount
		summary.TotalRepaidAmount += loan.PaidAmount
		summary.TotalOutstanding += loan.RepayAmount - loan.PaidAmount

		switch loan.Status {
		case "ACTIVE":
			summary.TotalActiveLoans++
		case "COMPLETED":
			summary.TotalCompletedLoans++
		case "DEFAULTED":
			summary.TotalDefaultedLoans++
		}

		branchLoanCount[loan.BranchName]++
		officerLoanCount[loan.LoanOfficer]++
	}

	var maxBranchLoans, maxOfficerLoans int64
	for branch, count := range branchLoanCount {
		if count > maxBranchLoans {
			maxBranchLoans = count
			summary.MostIssuedLoanBranch = branch
		}
	}
	for officer, count := range officerLoanCount {
		if count > maxOfficerLoans {
			maxOfficerLoans = count
			summary.MostLoansOfficer = officer
		}
	}

	return rslt, summary, nil
}



// func (r *LoanRepository) GetReportLoanData(ctx context.Context, filters services.ReportFilters) ([]services.LoanReportData, services.LoanSummary, error) {
// 	loans, err := r.GetLoansReportData(ctx, GetLoansReportDataParams{
// 		StartDate: filters.StartDate,
// 		EndDate: filters.EndDate,
// 	})
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, services.LoanSummary{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loans found")
// 		}

// 		return nil, services.LoanSummary{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get report loan data: %s", err.Error())
// 	}

// 	rslt := make([]services.LoanReportData, len(loans))

// 	summary := services.LoanSummary{}
// 	for i, loan := range loans {
// 		dueDate := ""
// 		disbursedDate := ""
// 		if loan.DueDate.Valid {
// 			dueDate = loan.DueDate.Time.Format("2006-02-01")
// 		}
// 		if loan.DisbursedDate.Valid {
// 			disbursedDate = loan.DueDate.Time.Format("2006-02-01")
// 		}
// 		rslt[i] = services.LoanReportData{
// 			LoanID: loan.LoanID,
// 			ClientName: loan.ClientName,
// 			BranchName: loan.BranchName,
// 			LoanOfficer: loan.LoanOfficer,
// 			LoanAmount: loan.LoanAmount,
// 			RepayAmount: loan.RepayAmount,
// 			PaidAmount: loan.PaidAmount,
// 			OutstandingAmount: pkg.InterfaceFloat64(loan.OutstandingAmount),
// 			Status: loan.Status,
// 			TotalInstallments: loan.TotalInstallments,
// 			PaidInstallments: loan.PaidInstallments,
// 			DueDate: dueDate,
// 			DisbursedDate: disbursedDate,
// 			DefaultRisk: pkg.InterfaceFloat64(loan.DefaultRisk),
// 		}
// 	}

// 	return rslt, summary, nil
// }

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

func (r *LoanRepository) GetLoanEvents(ctx context.Context) ([]repository.LoanEvent, error) {
	events, err := r.GetLoanEventsHelper(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return []repository.LoanEvent{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no loans found")
		}
		return []repository.LoanEvent{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get loan events: %s", err.Error())
	}

	rslt := []repository.LoanEvent{}
	for _, event := range events {
		var disbursedDate *string
		if event.DisbursedDate.Valid {
			d := event.DisbursedDate.Time.Format("2006-01-02")
			disbursedDate = &d
		}

		var dueDate *string
		if event.DueDate.Valid{
			d := event.DueDate.Time.Format("2006-01-02")
			dueDate = &d
		}

		var paymentDue *float64
		if event.PaymentDue.Valid {
			paymentDue = &event.PaymentDue.Float64
		}

		rslt = append(rslt, repository.LoanEvent{
			ID: fmt.Sprintf("LN%03d", event.LoanID),
			LoanID:      event.LoanID,
			ClientName:  event.ClientName,
			LoanAmount:  event.LoanAmount,
			Date: disbursedDate,
			Type:        "disbursed",
			Title: "Loan Disbursed",
			AllDay: false,
		})

		if dueDate != nil {
			rslt = append(rslt, repository.LoanEvent{
				ID: fmt.Sprintf("LN%03d", event.LoanID),
				LoanID:      event.LoanID,
				ClientName:  event.ClientName,
				LoanAmount:  event.LoanAmount,
				Date:     dueDate,
				PaymentDue:  paymentDue,
				Type:        "due",
				Title: "Payment Due",
				AllDay: false,
			})
		}
	}

	return rslt, nil
}

func convertGeneratedInstallment(installment generated.Installment) repository.Installment {
	return repository.Installment{
		ID:              installment.ID,
		LoanID:          installment.LoanID,
		InstallmentNo:   installment.InstallmentNumber,
		Amount:          installment.AmountDue,
		RemainingAmount: installment.RemainingAmount,
		Paid:            installment.Paid,
		PaidAt:          installment.PaidAt.Time.Format("2006-02-01"),
		DueDate:         installment.DueDate.Format("2006-02-01"),
	}
}

func convertGeneratedInstallmentList(installments []generated.Installment) []repository.Installment {
	rslt := make([]repository.Installment, len(installments))
	for idx, installment := range installments {
		rslt[idx] = repository.Installment{
			ID:              installment.ID,
			LoanID:          installment.LoanID,
			InstallmentNo:   installment.InstallmentNumber,
			Amount:          installment.AmountDue,
			RemainingAmount: installment.RemainingAmount,
			Paid:            installment.Paid,
			PaidAt:          installment.PaidAt.Time.Format("2006-02-01"),
			DueDate:         installment.DueDate.Format("2006-02-01"),
		}
	}
	
	return rslt
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
