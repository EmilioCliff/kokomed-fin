package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.ProductRepository = (*ProductRepository)(nil)

type ProductRepository struct {
	db      *Store
	queries generated.Querier
}

func NewProductRepository(db *Store) *ProductRepository {
	return &ProductRepository{
		db:      db,
		queries: generated.New(db.db),
	}
}

func (r *ProductRepository) GetAllProducts(ctx context.Context, pgData *pkg.PaginationMetadata) ([]repository.Product, error) {
	products, err := r.queries.ListProducts(ctx, generated.ListProductsParams{
		Limit:  int32(pgData.PageSize),
		Offset: int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no products found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get products: %s", err.Error())
	}

	result := make([]repository.Product, len(products))

	for i, product := range products {
		result[i] = convertGeneratedProducts(product)
	}

	return result, nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id uint32) (repository.Product, error) {
	product, err := r.queries.GetProduct(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.Product{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no product found")
		}

		return repository.Product{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get product: %s", err.Error())
	}

	return convertGeneratedProducts(product), nil
}

func (r *ProductRepository) ListProductByBranch(ctx context.Context, branchID uint32, pgData *pkg.PaginationMetadata) ([]repository.Product, error) {
	products, err := r.queries.ListProductsByBranch(ctx, generated.ListProductsByBranchParams{
		BranchID: branchID,
		Limit:    int32(pgData.PageSize),
		Offset:   int32(pkg.CalculateOffset(pgData.CurrentPage, pgData.PageSize)),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no products found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get products: %s", err.Error())
	}

	result := make([]repository.Product, len(products))

	for i, product := range products {
		result[i] = convertGeneratedProducts(product)
	}

	return result, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *repository.Product) (repository.Product, error) {
	execRslt, err := r.queries.CreateProduct(ctx, generated.CreateProductParams{
		BranchID:       product.BranchID,
		LoanAmount:     product.LoanAmount,
		RepayAmount:    product.RepayAmount,
		InterestAmount: product.InterestAmount,
		UpdatedBy:      product.UpdatedBy,
	})
	if err != nil {
		return repository.Product{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create product: %s", err.Error())
	}

	id, err := execRslt.LastInsertId()
	if err != nil {
		return repository.Product{}, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get last insert id: %s", err.Error())
	}

	product.ID = uint32(id)

	return *product, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id uint32) error {
	err := r.queries.DeleteProduct(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return pkg.Errorf(pkg.NOT_FOUND_ERROR, "no product found")
		}

		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to delete product: %s", err.Error())
	}

	return nil
}

func convertGeneratedProducts(product generated.Product) repository.Product {
	return repository.Product{
		ID:             product.ID,
		BranchID:       product.BranchID,
		LoanAmount:     product.LoanAmount,
		RepayAmount:    product.RepayAmount,
		InterestAmount: product.InterestAmount,
		UpdatedBy:      product.UpdatedBy,
		UpdatedAt:      product.UpdatedAt,
		CreatedAt:      product.CreatedAt,
	}
}
