package application

import (
	"context"
	domaincatalog "github.com/mytheresa/go-hiring-challenge/app/catalog/domain"
	"github.com/mytheresa/go-hiring-challenge/app/variants/domain"
)

type GetProductByIDUseCase struct {
	productRepo domaincatalog.Repository
	variantRepo domain.VariantRepository
}

func NewGetProductByIDUseCase(vRepo domain.VariantRepository) *GetProductByIDUseCase {
	return &GetProductByIDUseCase{
		variantRepo: vRepo,
	}
}

// Execute retrieves a product by its ID along with its variants.
// It returns a ProductWithVariants and an error if the product is not found or any other error occurs.
func (uc *GetProductByIDUseCase) Execute(ctx context.Context, req domain.GetProductByIDRequest) (*domain.ProductWithVariants, error) {
	product, err := uc.variantRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}

	variants, err := uc.variantRepo.FindVariantsByProductID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &domain.ProductWithVariants{
		Product:  product.Product,
		Variants: variants,
	}, nil
}
