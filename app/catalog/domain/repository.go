package domain

import "context"

type GetCatalogUseCaseInterface interface {
	Execute(ctx context.Context, req GetCatalogRequest) ([]Product, error)
}

type GetCatalogRequest struct {
	Offset   int
	Limit    int
	Category string
	PriceLt  *float64
}

type Repository interface {
	// FindAll retrieves products based on the provided criteria.
	// It returns a slice of Product and an error if any.
	FindAll(ctx context.Context, offset int, limit int, category string, priceLt *float64) ([]Product, error)

	// FindByID retrieves a product by its ID.
	// It returns a pointer to Product and an error if the product is not found or any other error occurs.
	FindByID(ctx context.Context, id uint) (*Product, error)
}
