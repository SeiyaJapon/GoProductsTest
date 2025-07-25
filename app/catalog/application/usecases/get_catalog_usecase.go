package usecases

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/catalog/domain/product"
)

type GetCatalogRequest struct {
	Offset   int
	Limit    int
	Category string
	PriceLt  *float64
}

type GetCatalogUseCaseInterface interface {
	Execute(ctx context.Context, req GetCatalogRequest) ([]product.Product, error)
}

type GetCatalogUseCase struct {
	repo product.Repository
}

func NewGetCatalogUseCase(repo product.Repository) *GetCatalogUseCase {
	return &GetCatalogUseCase{repo: repo}
}

func (uc *GetCatalogUseCase) Execute(ctx context.Context, req GetCatalogRequest) ([]product.Product, error) {
	return uc.repo.FindAll(ctx, req.Offset, req.Limit, req.Category, req.PriceLt)
}
