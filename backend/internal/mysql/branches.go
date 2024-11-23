package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
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

func (r *BranchRepository) CreateBranch(ctx context.Context, branch *repository.Branch) (repository.Branch, error) {
	execResult, err := r.queries.CreateBranch(ctx, branch.Name)
	if err != nil {
		return repository.Branch{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create branch: %s", err.Error())
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.Branch{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	branch.ID = uint32(id)

	return *branch, nil
}

func (r *BranchRepository) ListBranches(ctx context.Context) ([]repository.Branch, error) {
	branches, err := r.queries.ListBranches(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no branches found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get branches: %s", err.Error())
	}

	result := make([]repository.Branch, len(branches))

	for idx, branch := range branches {
		result[idx] = repository.Branch{
			ID:   branch.ID,
			Name: branch.Name,
		}
	}

	return result, nil
}

func (r *BranchRepository) GetBranchByID(ctx context.Context, id uint32) (repository.Branch, error) {
	branch, err := r.queries.GetBranch(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.Branch{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no branch found")
		}

		return repository.Branch{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get branch: %s", err.Error())
	}

	return repository.Branch{
		ID:   branch.ID,
		Name: branch.Name,
	}, nil
}

func (r *BranchRepository) UpdateBranch(ctx context.Context, name string, id uint32) (repository.Branch, error) {
	_, err := r.queries.UpdateBranch(ctx, generated.UpdateBranchParams{
		Name: name,
		ID:   id,
	})
	if err != nil {
		return repository.Branch{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update branch: %s", err.Error())
	}

	return repository.Branch{
		ID:   id,
		Name: name,
	}, nil
}
