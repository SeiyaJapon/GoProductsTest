package product

import "context"

type Repository interface {
	FindAll(ctx context.Context, offset int, limit int, category string, priceLt *float64) ([]Product, error)
	FindByID(ctx context.Context, id uint) (*Product, error)
}
