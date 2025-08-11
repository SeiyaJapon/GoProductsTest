package catalog

import (
	"context"
)

type GetCatalogRequest struct {
	Offset   int
	Limit    int
	Category string
	PriceLt  *float64
}

type GetCatalogUseCaseInterface interface {
	Execute(ctx context.Context, req GetCatalogRequest) ([]Product, error)
}

type Repository interface {
	FindAll(ctx context.Context, offset int, limit int, category string, priceLt *float64) ([]Product, error)
	FindByID(ctx context.Context, id uint) (*Product, error)
}

type GetCatalogUseCase struct {
	repo Repository
}

func NewGetCatalogUseCase(repo Repository) *GetCatalogUseCase {
	return &GetCatalogUseCase{repo: repo}
}

func (uc *GetCatalogUseCase) Execute(ctx context.Context, req GetCatalogRequest) ([]Product, error) {
	return uc.repo.FindAll(ctx, req.Offset, req.Limit, req.Category, req.PriceLt)
}
