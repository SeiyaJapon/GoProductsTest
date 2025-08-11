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
	FindAll(ctx context.Context, offset int, limit int, category string, priceLt *float64) ([]Product, error)
	FindByID(ctx context.Context, id uint) (*Product, error)
}
