package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

func (r *LoanRepository) CreateLoan(ctx context.Context, loan *repository.Loan) (repository.Loan, error) {
	if loan.DueDate != nil && loan.DisbursedBy != nil && loan.DisbursedOn != nil {
		hasActiveLoan, err := r.queries.CheckActiveLoanForClient(ctx, loan.ClientID)
		if err != nil {
			return repository.Loan{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to check if client has an active loan: %s", err.Error())
		}

		if hasActiveLoan {
			return repository.Loan{}, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "client already has an active loan")
		}
	}

	err := r.db.ExecTx(ctx, func(q *generated.Queries) error {
		// create the loan
		params := generated.CreateLoanParams{
			ProductID:          loan.ProductID,
			ClientID:           loan.ClientID,
			LoanOfficer:        loan.LoanOfficerID,
			ApprovedBy:         loan.ApprovedBy,
			TotalInstallments:  loan.TotalInstallments,
			InstallmentsPeriod: loan.InstallmentsPeriod,
			ProcessingFee:      loan.ProcessingFee,
			FeePaid:            loan.FeePaid,
			CreatedBy:          loan.CreatedBy,
			Status:             generated.LoansStatusINACTIVE,
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
			params.Status = generated.LoansStatusACTIVE
		}

		execResult, err := q.CreateLoan(ctx, params)
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create loan: %s", err.Error())
		}

		id, err := execResult.LastInsertId()
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
		}

		loan.ID = uint32(id)

		// create loan installments if loan is disbursed already(loan is ACTIVE)
		if params.Status == generated.LoansStatusACTIVE {
			if err := helperCreateInstallation(ctx, q, loan.ID, loan.ProductID, params.TotalInstallments, params.InstallmentsPeriod); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return repository.Loan{}, err
	}

	return *loan, nil
}

func (r *LoanRepository) DisburseLoan(ctx context.Context, disburseLoan *repository.DisburseLoan) error {
	loan, err := r.GetLoanByID(ctx, disburseLoan.ID)
	if err != nil {
		return err
	}

	hasActiveLoan, err := r.queries.CheckActiveLoanForClient(ctx, loan.ClientID)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to check if client has an active loan: %s", err.Error())
	}

	if hasActiveLoan {
		return pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "client already has an active loan")
	}

	err = r.db.ExecTx(ctx, func(q *generated.Queries) error {
		_, err = r.queries.DisburseLoan(ctx, generated.DisburseLoanParams{
			ID:     disburseLoan.ID,
			Status: generated.LoansStatusACTIVE,
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
				Time:  disburseLoan.DueDate.AddDate(0, 0, int(loan.InstallmentsPeriod)*int(loan.TotalInstallments)),
			},
		})
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to disburse loan: %s", err.Error())
		}

		// create installments
		if err = helperCreateInstallation(ctx, q, loan.ID, loan.ProductID, loan.TotalInstallments, loan.InstallmentsPeriod); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func helperCreateInstallation(ctx context.Context, q *generated.Queries, loanID, productID, totalInstallment, intallmentPeriod uint32) error {
	repayAmout, err := q.GetProductRepayAmount(ctx, productID)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get product repay amount: %s", err.Error())
	}

	installmentAmount := repayAmout / float64(totalInstallment)
	firstDueDate := time.Now().AddDate(0, 0, int(intallmentPeriod))

	for i := 0; i < int(totalInstallment); i++ {
		dueDate := firstDueDate.AddDate(0, 0, i*int(intallmentPeriod))

		_, err := q.CreateInstallment(ctx, generated.CreateInstallmentParams{
			LoanID:            loanID,
			InstallmentNumber: uint32(i + 1),
			AmountDue:         installmentAmount,
			RemainingAmount:   installmentAmount,
			DueDate:           dueDate,
		})
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create installment: %s", err.Error())
		}
	}

	return nil
}
