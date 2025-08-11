package domain

import (
	"context"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/domain"
)

type GetProductByIDRequest struct {
	ID uint
}

type ProductWithVariants struct {
	Product  *domain.Product
	Variants []Variant
}

type GetProductByIDUseCaseInterface interface {
	// Execute retrieves a product by its ID along with its variants.
	// It returns a ProductWithVariants and an error if any.
	Execute(ctx context.Context, req GetProductByIDRequest) (*ProductWithVariants, error)
}

type VariantRepository interface {
	// FindByID retrieves a product by its ID.
	// It returns a ProductWithVariants and an error if the product is not found or any other error occurs.
	FindByID(ctx context.Context, id uint) (*ProductWithVariants, error)

	// FindVariantsByProductID retrieves all variants for a given product ID.
	// It returns a slice of Variant and an error if any.
	FindVariantsByProductID(ctx context.Context, productID uint) ([]Variant, error)
}
