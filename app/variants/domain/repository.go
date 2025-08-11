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
	Execute(ctx context.Context, req GetProductByIDRequest) (*ProductWithVariants, error)
}

type VariantRepository interface {
	FindByID(ctx context.Context, id uint) (*ProductWithVariants, error)
	FindVariantsByProductID(ctx context.Context, productID uint) ([]Variant, error)
}
