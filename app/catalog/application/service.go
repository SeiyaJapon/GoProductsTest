package application

import (
	"context"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/domain"
)

type GetCatalogUseCase struct {
	repo domain.Repository
}

func NewGetCatalogUseCase(repo domain.Repository) *GetCatalogUseCase {
	return &GetCatalogUseCase{repo: repo}
}

func (uc *GetCatalogUseCase) Execute(ctx context.Context, req domain.GetCatalogRequest) ([]domain.Product, error) {
	return uc.repo.FindAll(ctx, req.Offset, req.Limit, req.Category, req.PriceLt)
}
