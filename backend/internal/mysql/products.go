package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
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

func (r *ProductRepository) GetAllProducts(
	ctx context.Context,
	search *string,
	pgData *pkg.PaginationMetadata,
) ([]repository.Product, pkg.PaginationMetadata, error) {
	params := generated.ListProductsByCategoryParams{
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

	products, err := r.queries.ListProductsByCategory(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.PaginationMetadata{}, pkg.Errorf(
				pkg.NOT_FOUND_ERROR,
				"no products found",
			)
		}

		return nil, pkg.PaginationMetadata{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get products: %s",
			err.Error(),
		)
	}

	totalProducts, err := r.queries.CountLoansByCategory(ctx, params2)
	if err != nil {
		return nil, pkg.PaginationMetadata{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get total products: %s",
			err.Error(),
		)
	}

	result := make([]repository.Product, len(products))
	for i, product := range products {
		result[i] = repository.Product{
			ID:             product.ID,
			BranchID:       product.BranchID,
			BranchName:     &product.BranchName,
			LoanAmount:     product.LoanAmount,
			RepayAmount:    product.RepayAmount,
			InterestAmount: product.InterestAmount,
			UpdatedBy:      product.UpdatedBy,
			UpdatedAt:      product.UpdatedAt,
			CreatedAt:      product.CreatedAt,
		}
	}

	return result, pkg.CreatePaginationMetadata(
		uint32(totalProducts),
		pgData.PageSize,
		pgData.CurrentPage,
	), nil
}

func (r *ProductRepository) GetProductByID(
	ctx context.Context,
	id uint32,
) (repository.Product, error) {
	product, err := r.queries.GetProduct(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.Product{}, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no product found")
		}

		return repository.Product{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get product: %s",
			err.Error(),
		)
	}

	return repository.Product{
		ID:             product.ID,
		BranchID:       product.BranchID,
		LoanAmount:     product.LoanAmount,
		RepayAmount:    product.RepayAmount,
		InterestAmount: product.InterestAmount,
		UpdatedBy:      product.UpdatedBy,
		UpdatedAt:      product.UpdatedAt,
		CreatedAt:      product.CreatedAt,
		BranchName:     &product.BranchName,
	}, nil
}

func (r *ProductRepository) ListProductByBranch(
	ctx context.Context,
	branchID uint32,
	pgData *pkg.PaginationMetadata,
) ([]repository.Product, error) {
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

func (r *ProductRepository) CreateProduct(
	ctx context.Context,
	product *repository.Product,
) (repository.Product, error) {
	execRslt, err := r.queries.CreateProduct(ctx, generated.CreateProductParams{
		BranchID:       product.BranchID,
		LoanAmount:     product.LoanAmount,
		RepayAmount:    product.RepayAmount,
		InterestAmount: product.InterestAmount,
		UpdatedBy:      product.UpdatedBy,
	})
	if err != nil {
		return repository.Product{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to create product: %s",
			err.Error(),
		)
	}

	id, err := execRslt.LastInsertId()
	if err != nil {
		return repository.Product{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get last insert id: %s",
			err.Error(),
		)
	}

	branch, err := r.queries.GetBranch(ctx, product.BranchID)
	if err != nil {
		return repository.Product{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get created product branch: %s",
			err.Error(),
		)
	}

	product.ID = uint32(id)
	product.BranchName = &branch.Name

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

func (r *ProductRepository) GetReportProductData(
	ctx context.Context,
	filters services.ReportFilters,
) ([]services.ProductReportData, services.ProductSummary, error) {
	products, err := r.GetProductReportData(ctx, GetProductReportDataParams{
		StartDate: filters.StartDate,
		EndDate:   filters.EndDate,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, services.ProductSummary{}, pkg.Errorf(
				pkg.NOT_FOUND_ERROR,
				"no product found",
			)
		}

		return nil, services.ProductSummary{}, pkg.Errorf(
			pkg.INTERNAL_ERROR,
			"failed to get report product data: %s",
			err.Error(),
		)
	}

	rslt := make([]services.ProductReportData, len(products))

	var totalActiveLoanAmount int64
	var mostPopularProduct string
	var maxLoans int64
	for i, product := range products {
		totalActiveLoanAmount += product.ActiveLoans

		if product.TotalLoansIssued > maxLoans {
			maxLoans = product.TotalLoansIssued
			mostPopularProduct = product.ProductName
		}

		rslt[i] = services.ProductReportData{
			ProductName:       product.ProductName,
			LoansIssued:       product.TotalLoansIssued,
			ActiveLoans:       product.ActiveLoans,
			CompletedLoans:    product.CompletedLoans,
			DefaultedLoans:    product.DefaultedLoans,
			AmountDisbursed:   pkg.InterfaceFloat64(product.TotalAmountDisbursed),
			AmountRepaid:      pkg.InterfaceFloat64(product.TotalAmountRepaid),
			OutstandingAmount: pkg.InterfaceFloat64(product.TotalOutstandingAmount),
			DefaultRate:       pkg.InterfaceFloat64(product.DefaultRate),
		}
	}

	summary := services.ProductSummary{
		TotalProducts:         int64(len(products)),
		TotalActiveLoanAmount: totalActiveLoanAmount,
		MostPopularProduct:    mostPopularProduct,
		MaxLoans:              maxLoans,
	}

	return rslt, summary, nil
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
