package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

func (r *LoanRepository) CreateLoan(
	ctx context.Context,
	loan *repository.Loan,
) (repository.LoanFullData, error) {
	// check if client has alread an active loan
	if loan.DueDate != nil && loan.DisbursedBy != nil && loan.DisbursedOn != nil {
		hasActiveLoan, err := r.queries.CheckActiveLoanForClient(ctx, loan.ClientID)
		if err != nil {
			return repository.LoanFullData{}, pkg.Errorf(
				pkg.INTERNAL_ERROR,
				"failed to check if client has an active loan: %s",
				err.Error(),
			)
		}

		if hasActiveLoan {
			return repository.LoanFullData{}, pkg.Errorf(
				pkg.ALREADY_EXISTS_ERROR,
				"client already has an active loan",
			)
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
			if err := helperCreateInstallation(ctx, q, *loan.DisbursedOn, loan.ID, loan.ProductID, params.TotalInstallments, params.InstallmentsPeriod); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return repository.LoanFullData{}, err
	}

	updateLoan, err := r.queries.GetLoanFullData(ctx, loan.ID)
	if err != nil {
		return repository.LoanFullData{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get created loan: %s",
			err.Error(),
		)
	}

	return convertGetLoanFullDataRowToRepo(&updateLoan), nil
}

func (r *LoanRepository) DisburseLoan(
	ctx context.Context,
	disburseLoan *repository.DisburseLoan,
) (uint32, error) {
	loan, err := r.GetLoanByID(ctx, disburseLoan.ID)
	if err != nil {
		return 0, err
	}

	hasActiveLoan, err := r.queries.CheckActiveLoanForClient(ctx, loan.ClientID)
	if err != nil {
		return 0, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to check if client has an active loan: %s",
			err.Error(),
		)
	}

	if disburseLoan.DisbursedOn != nil {
		if hasActiveLoan {
			return 0, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "client already has an active loan")
		}
	}

	err = r.db.ExecTx(ctx, func(q *generated.Queries) error {
		params := generated.DisburseLoanParams{
			ID: disburseLoan.ID,
			DisbursedBy: sql.NullInt32{
				Valid: true,
				Int32: int32(disburseLoan.DisbursedBy),
			},
		}

		if disburseLoan.DisbursedOn != nil {
			params.DisbursedOn = sql.NullTime{
				Valid: true,
				Time:  *disburseLoan.DisbursedOn,
			}

			dueDate := (*disburseLoan.DisbursedOn).AddDate(
				0,
				0,
				int(loan.InstallmentsPeriod)*int(loan.TotalInstallments),
			)

			params.DueDate = sql.NullTime{
				Valid: true,
				Time:  dueDate,
			}
		}

		if disburseLoan.Status != nil {
			params.Status = generated.NullLoansStatus{
				Valid:       true,
				LoansStatus: generated.LoansStatus(*disburseLoan.Status),
			}
		}

		if disburseLoan.FeePaid != nil {
			params.FeePaid = sql.NullBool{
				Valid: true,
				Bool:  true,
			}
		}

		_, err = r.queries.DisburseLoan(ctx, params)
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to disburse loan: %s", err.Error())
		}

		// here we will create installments if and only if the status was changed to active
		if disburseLoan.Status != nil &&
			generated.LoansStatus(*disburseLoan.Status) == generated.LoansStatusACTIVE {
			if err = helperCreateInstallation(ctx, q, *disburseLoan.DisbursedOn, loan.ID, loan.ProductID, loan.TotalInstallments, loan.InstallmentsPeriod); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return loan.ClientID, nil
}

func helperCreateInstallation(
	ctx context.Context,
	q *generated.Queries,
	disbursedDate time.Time,
	loanID, productID, totalInstallment, intallmentPeriod uint32,
) error {
	repayAmout, err := q.GetProductRepayAmount(ctx, productID)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get product repay amount: %s", err.Error())
	}

	installmentAmount := repayAmout / float64(totalInstallment)
	firstDueDate := disbursedDate.AddDate(0, 0, int(intallmentPeriod))

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

func convertGetLoanFullDataRowToRepo(loan *generated.GetLoanFullDataRow) repository.LoanFullData {
	rsp := repository.LoanFullData{
		ID: loan.ID,
		Product: repository.ProductShort{
			ID:             loan.ProductID,
			BranchName:     loan.ProductBranchName, // You might need to join the product's branch name if required
			LoanAmount:     loan.LoanAmount,
			RepayAmount:    loan.RepayAmount,
			InterestAmount: loan.InterestAmount,
		},
		Client: repository.ClientShort{
			ID:          loan.ClientID,
			FullName:    loan.ClientName,
			PhoneNumber: loan.ClientPhone,
			Active:      loan.ClientActive,
			BranchName:  loan.ClientBranchName,
		},
		LoanOfficer: repository.UserShortResponse{
			ID:          loan.LoanOfficer,
			FullName:    loan.LoanOfficerName,
			Email:       loan.LoanOfficerEmail,
			PhoneNumber: loan.LoanOfficerPhone,
		},
		LoanPurpose: pkg.StringPtr(""),
		DueDate:     &time.Time{},
		ApprovedBy: repository.UserShortResponse{
			ID:          loan.ApprovedBy,
			FullName:    loan.ApprovedByName,
			Email:       loan.ApprovedByEmail,
			PhoneNumber: loan.ApprovedByPhone,
		},
		DisbursedOn:        &time.Time{},
		TotalInstallments:  loan.TotalInstallments,
		InstallmentsPeriod: loan.InstallmentsPeriod,
		Status:             string(loan.Status),
		ProcessingFee:      loan.ProcessingFee,
		FeePaid:            loan.FeePaid,
		PaidAmount:         loan.PaidAmount,
		RemainingAmount:    loan.RepayAmount - loan.PaidAmount,
		CreatedBy: repository.UserShortResponse{
			ID:          loan.CreatedBy,
			FullName:    loan.CreatedByName.String,
			Email:       loan.CreatedByEmail.String,
			PhoneNumber: loan.CreatedByPhone.String,
		},
		CreatedAt: loan.CreatedAt,
	}

	if loan.DueDate.Valid {
		rsp.DueDate = &loan.DueDate.Time
	}

	if loan.DisbursedOn.Valid {
		rsp.DisbursedOn = &loan.DisbursedOn.Time
	}

	if loan.LoanPurpose.Valid {
		rsp.LoanPurpose = &loan.LoanPurpose.String
	}

	if loan.UpdatedBy.Valid {
		rsp.UpdatedBy = repository.UserShortResponse{
			ID:          uint32(loan.UpdatedBy.Int32),
			FullName:    loan.UpdatedByName.String,
			Email:       loan.UpdatedByEmail.String,
			PhoneNumber: loan.UpdatedByPhone.String,
		}
	}

	if loan.DisbursedBy.Valid {
		rsp.DisbursedBy = repository.UserShortResponse{
			ID:          uint32(loan.DisbursedBy.Int32),
			FullName:    loan.DisbursedByName.String,
			Email:       loan.DisbursedByEmail.String,
			PhoneNumber: loan.DisbursedByPhone.String,
		}
	}

	return rsp
}
