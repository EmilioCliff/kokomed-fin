package mock

import (
	"context"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

var _ repository.ProductRepository = (*MockProductRepository)(nil)

type MockProductRepository struct {
	mockGetAllProductsFunc      func(ctx context.Context, pgData *pkg.PaginationMetadata) ([]repository.Product, error)
	mockGetProductByIDFunc      func(ctx context.Context, id uint32) (repository.Product, error)
	mockListProductByBranchFunc func(ctx context.Context, branchID uint32, pgData *pkg.PaginationMetadata) ([]repository.Product, error)
	mockCreateProductFunc       func(ctx context.Context, product *repository.Product) (repository.Product, error)
	mockUpdateProductFunc       func(ctx context.Context, product *repository.UpdateProduct) (repository.Product, error)
	mockDeleteProductFunc       func(ctx context.Context, id uint32) error
}

func (m *MockProductRepository) GetAllProducts(ctx context.Context, pgData *pkg.PaginationMetadata) ([]repository.Product, error) {
	return m.mockGetAllProductsFunc(ctx, pgData)
}

func (m *MockProductRepository) GetProductByID(ctx context.Context, id uint32) (repository.Product, error) {
	return m.mockGetProductByIDFunc(ctx, id)
}

func (m *MockProductRepository) ListProductByBranch(
	ctx context.Context,
	branchID uint32,
	pgData *pkg.PaginationMetadata,
) ([]repository.Product, error) {
	return m.mockListProductByBranchFunc(ctx, branchID, pgData)
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, product *repository.Product) (repository.Product, error) {
	return m.mockCreateProductFunc(ctx, product)
}

func (m *MockProductRepository) UpdateProduct(ctx context.Context, product *repository.UpdateProduct) (repository.Product, error) {
	return m.mockUpdateProductFunc(ctx, product)
}

func (m *MockProductRepository) DeleteProduct(ctx context.Context, id uint32) error {
	return m.mockDeleteProductFunc(ctx, id)
}
