package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.NonPostedRepository = (*NonPostedRepository)(nil)

type NonPostedRepository struct {
	db      *Store
	queries generated.Querier
}

func NewNonPostedRepository(db *Store) *NonPostedRepository {
	return &NonPostedRepository{
		db:      db,
		queries: generated.New(db.db),
	}
}

func (r *NonPostedRepository) CreateNonPosted(ctx context.Context, nonPosted *repository.NonPosted) (repository.NonPosted, error) {
	execResult, err := r.queries.CreateNonPosted(ctx, generated.CreateNonPostedParams{
		TransactionNumber: nonPosted.TransactionNumber,
		AccountNumber:     nonPosted.AccountNumber,
		PhoneNumber:       nonPosted.PhoneNumber,
		PayingName:        nonPosted.PayingName,
		Amount:            nonPosted.Amount,
		PaidDate:          nonPosted.PaidDate,
	})
	if err != nil {
		return repository.NonPosted{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create non posted: %s", err.Error())
	}

	id, err := execResult.LastInsertId()
	if err != nil {
		return repository.NonPosted{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	nonPosted.ID = uint32(id)

	return *nonPosted, nil
}

func (r *NonPostedRepository) GetNonPosted(ctx context.Context, id uint32) (repository.NonPosted, error) {
	nonPosted, err := r.queries.GetNonPosted(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.NonPosted{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no non posted found")
		}

		return repository.NonPosted{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get non posted: %s", err.Error())
	}

	return convertGenerateNonPosted(nonPosted), nil
}

func (r *NonPostedRepository) AssignNonPosted(ctx context.Context, id uint32, assignedTo uint32) (repository.NonPosted, error) {
	execResult, err := r.queries.AssignNonPosted(ctx, generated.AssignNonPostedParams{
		ID: id,
		AssignTo: sql.NullInt32{
			Valid: true,
			Int32: int32(assignedTo),
		},
	})
	if err != nil {
		return repository.NonPosted{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to assign non posted: %s", err.Error())
	}

	genId, err := execResult.LastInsertId()
	if err != nil {
		return repository.NonPosted{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	return r.GetNonPosted(ctx, uint32(genId))
}

func (r *NonPostedRepository) ListNonPosted(ctx context.Context, pgData *pkg.PaginationMetadata) ([]repository.NonPosted, error) {
	nonPosteds, err := r.queries.ListAllNonPosted(ctx, generated.ListAllNonPostedParams{
		Limit:  pkg.GetPageSize(),
		Offset: pkg.CalculateOffset(pgData.CurrentPage),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no non posted found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list non posted: %s", err.Error())
	}

	rslt := make([]repository.NonPosted, len(nonPosteds))

	for i, nonPosted := range nonPosteds {
		rslt[i] = convertGenerateNonPosted(nonPosted)
	}

	return rslt, nil
}

func (r *NonPostedRepository) ListUnassignedNonPosted(ctx context.Context, pgData *pkg.PaginationMetadata) ([]repository.NonPosted, error) {
	nonPosteds, err := r.queries.ListUnassignedNonPosted(ctx, generated.ListUnassignedNonPostedParams{
		Limit:  pkg.GetPageSize(),
		Offset: pkg.CalculateOffset(pgData.CurrentPage),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no non posted found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to list non posted: %s", err.Error())
	}

	rslt := make([]repository.NonPosted, len(nonPosteds))

	for i, nonPosted := range nonPosteds {
		rslt[i] = convertGenerateNonPosted(nonPosted)
	}

	return rslt, nil
}

func (r *NonPostedRepository) DeleteNonPosted(ctx context.Context, id uint32) error {
	err := r.queries.DeleteNonPosted(ctx, id)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to delete non posted: %s", err.Error())
	}

	return nil
}

func convertGenerateNonPosted(nonPosted generated.NonPosted) repository.NonPosted {
	var assignedTo *uint32

	if nonPosted.AssignTo.Valid {
		value := uint32(nonPosted.AssignTo.Int32)
		assignedTo = &value
	}

	return repository.NonPosted{
		ID:                nonPosted.ID,
		TransactionNumber: nonPosted.TransactionNumber,
		AccountNumber:     nonPosted.AccountNumber,
		PhoneNumber:       nonPosted.PhoneNumber,
		PayingName:        nonPosted.PayingName,
		Amount:            nonPosted.Amount,
		PaidDate:          nonPosted.PaidDate,
		AssignedTo:        assignedTo,
	}
}