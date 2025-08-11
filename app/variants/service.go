package variants

import (
	"context"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
)

type GetProductByIDRequest struct {
	ID uint
}

type GetProductByIDUseCase struct {
	productRepo catalog.Repository
	variantRepo VariantRepository
}

type ProductWithVariants struct {
	Product  *catalog.Product
	Variants []Variant
}

type GetProductByIDUseCaseInterface interface {
	Execute(ctx context.Context, req GetProductByIDRequest) (*ProductWithVariants, error)
}

type VariantRepository interface {
	FindByID(ctx context.Context, id uint) (*ProductWithVariants, error)
	FindVariantsByProductID(ctx context.Context, productID uint) ([]Variant, error)
}

func NewGetProductByIDUseCase(vRepo VariantRepository) *GetProductByIDUseCase {
	return &GetProductByIDUseCase{
		variantRepo: vRepo,
	}
}

func (uc *GetProductByIDUseCase) Execute(ctx context.Context, req GetProductByIDRequest) (*ProductWithVariants, error) {
	product, err := uc.variantRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	variants, err := uc.variantRepo.FindVariantsByProductID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &ProductWithVariants{
		Product:  product.Product,
		Variants: variants,
	}, nil
}
