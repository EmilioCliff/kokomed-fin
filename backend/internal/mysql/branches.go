package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.BranchRepository = (*BranchRepository)(nil)

type BranchRepository struct {
	db      *Store
	queries generated.Querier
}

func NewBranchRepository(db *Store) *BranchRepository {
	return &BranchRepository{
		db:      db,
		queries: generated.New(db.db),
	}
}

func (r *BranchRepository) CreateBranch(
	ctx context.Context,
	branch *repository.Branch,
) (repository.Branch, error) {
	execResult, err := r.queries.CreateBranch(ctx, branch.Name)
	if err != nil {
		return repository.Branch{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to create branch: %s",
			err.Error(),
		)
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.Branch{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get last insert id: %s",
			err.Error(),
		)
	}

	branch.ID = uint32(id)

	return *branch, nil
}

func (r *BranchRepository) ListBranches(
	ctx context.Context,
	search *string,
	pgData *pkg.PaginationMetadata,
) ([]repository.Branch, pkg.PaginationMetadata, error) {
	params := generated.ListBrachesByCategoryParams{
		Column1: "",
		Name:    "",
		Limit:   int32(pgData.PageSize),
		Offset:  int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
	}

	params2 := generated.CountLoansByCategoryParams{
		Column1: "",
		Name:    "",
	}

	if search != nil {
		searchValue := "%" + *search + "%"
		params.Column1 = "has_search"
		params.Name = searchValue

		params2.Column1 = "has_search"
		params2.Name = searchValue
	}

	branches, err := r.queries.ListBrachesByCategory(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.PaginationMetadata{}, pkg.Errorf(
				pkg.NOT_FOUND_ERROR,
				"no branches found",
			)
		}

		return nil, pkg.PaginationMetadata{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get branches: %s",
			err.Error(),
		)
	}

	totalBranches, err := r.queries.CountBranchesByCategory(
		ctx,
		generated.CountBranchesByCategoryParams(params2),
	)
	if err != nil {
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get total branches: %s",
			err.Error(),
		)
	}

	result := make([]repository.Branch, len(branches))

	for idx, branch := range branches {
		result[idx] = repository.Branch{
			ID:   branch.ID,
			Name: branch.Name,
		}
	}

	return result, pkg.CreatePaginationMetadata(
		uint32(totalBranches),
		pgData.PageSize,
		pgData.CurrentPage,
	), nil
}

func (r *BranchRepository) GetBranchByID(
	ctx context.Context,
	id uint32,
) (repository.Branch, error) {
	branch, err := r.queries.GetBranch(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.Branch{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no branch found")
		}

		return repository.Branch{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get branch: %s",
			err.Error(),
		)
	}

	return repository.Branch{
		ID:   branch.ID,
		Name: branch.Name,
	}, nil
}

func (r *BranchRepository) GetReportBranchData(
	ctx context.Context,
	filters services.ReportFilters,
) ([]services.BranchReportData, services.BranchSummary, error) {
	branches, err := r.queries.GetBranchReportData(ctx, generated.GetBranchReportDataParams{
		StartDate: sql.NullTime{
			Valid: true,
			Time:  filters.StartDate,
		},
		EndDate: sql.NullTime{
			Valid: true,
			Time:  filters.EndDate,
		},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, services.BranchSummary{}, pkg.Errorf(
				pkg.NOT_FOUND_ERROR,
				"no branches data found",
			)
		}
		return nil, services.BranchSummary{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get report branches data: %s",
			err.Error(),
		)
	}

	rslt := make([]services.BranchReportData, len(branches))

	var totalBranches int64
	var totalClients int64
	var totalUsers int64
	var totalLoansIssued int64

	var highestPerformingBranch string
	var mostClientsBranch string
	var highestDisbursed float64
	var maxClients int64

	for i, branch := range branches {
		disbursedAmount := pkg.InterfaceFloat64(branch.TotalDisbursedAmount)
		totalBranches++
		totalClients += branch.TotalClients
		totalUsers += branch.TotalUsers
		totalLoansIssued += branch.TotalLoansIssued

		if disbursedAmount > highestDisbursed {
			highestDisbursed = disbursedAmount
			highestPerformingBranch = branch.BranchName
		}

		if branch.TotalClients > maxClients {
			maxClients = branch.TotalClients
			mostClientsBranch = branch.BranchName
		}

		rslt[i] = services.BranchReportData{
			BranchName:       branch.BranchName,
			TotalClients:     branch.TotalClients,
			TotalUsers:       branch.TotalUsers,
			LoansIssued:      branch.TotalLoansIssued,
			TotalDisbursed:   disbursedAmount,
			TotalCollected:   pkg.InterfaceFloat64(branch.TotalCollectedAmount),
			TotalOutstanding: pkg.InterfaceFloat64(branch.TotalOutstandingAmount),
			DefaultRate:      pkg.InterfaceFloat64(branch.DefaultRate),
		}
	}

	summary := services.BranchSummary{
		TotalBranches:           totalBranches,
		TotalClients:            totalClients,
		TotalUsers:              totalUsers,
		HighestPerformingBranch: highestPerformingBranch,
		MostClientsBranch:       mostClientsBranch,
	}

	return rslt, summary, nil
}

func (r *BranchRepository) UpdateBranch(
	ctx context.Context,
	name string,
	id uint32,
) (repository.Branch, error) {
	_, err := r.queries.UpdateBranch(ctx, generated.UpdateBranchParams{
		Name: name,
		ID:   id,
	})
	if err != nil {
		return repository.Branch{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to update branch: %s",
			err.Error(),
		)
	}

	return repository.Branch{
		ID:   id,
		Name: name,
	}, nil
}
