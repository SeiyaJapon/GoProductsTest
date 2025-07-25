package usecases

import (
	"context"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/domain/product"
)

type GetProductByIDRequest struct {
	ID uint
}

type GetProductByIDUseCaseInterface interface {
	Execute(ctx context.Context, req GetProductByIDRequest) (*product.Product, error)
}

type GetProductByIDUseCase struct {
	repo product.Repository
}

func NewGetProductByIDUseCase(repo product.Repository) *GetProductByIDUseCase {
	return &GetProductByIDUseCase{repo: repo}
}

func (uc *GetProductByIDUseCase) Execute(ctx context.Context, req GetProductByIDRequest) (*product.Product, error) {
	return uc.repo.FindByID(ctx, req.ID)
}
