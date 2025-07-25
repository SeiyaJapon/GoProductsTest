package usecases

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/catalog/domain/product"
)

type mockProductRepositoryById struct {
	MockFindAllProducts []product.Product
	MockFindAllError    error

	MockFindByIDProduct *product.Product
	MockFindByIDError   error
}

func (m *mockProductRepositoryById) FindAll(ctx context.Context, offset int, limit int, category string, priceLt *float64) ([]product.Product, error) {
	return m.MockFindAllProducts, m.MockFindAllError
}

func (m *mockProductRepositoryById) FindByID(ctx context.Context, id uint) (*product.Product, error) {
	return m.MockFindByIDProduct, m.MockFindByIDError
}

func TestGetProductByIDUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should return product when repository finds it", func(t *testing.T) {
		expectedProduct := &product.Product{
			ID:    1,
			Code:  "P001",
			Price: 75.0,
			Category: product.Category{
				ID:   strconv.Itoa(201),
				Code: "SHOES",
				Name: "Shoes",
			},
		}

		mockRepo := &mockProductRepositoryById{
			MockFindByIDProduct: expectedProduct,
			MockFindByIDError:   nil,
		}

		useCase := NewGetProductByIDUseCase(mockRepo)

		request := GetProductByIDRequest{ID: 1}

		foundProduct, err := useCase.Execute(ctx, request)

		if err != nil {
			t.Fatalf("expected no error, but got: %v", err)
		}
		if foundProduct == nil {
			t.Fatal("expected a product, got nil")
		}
		if foundProduct.ID != expectedProduct.ID || foundProduct.Code != expectedProduct.Code {
			t.Errorf("expected product %+v, got %+v", expectedProduct, foundProduct)
		}
	})

	t.Run("should return nil and no error when product is not found", func(t *testing.T) {
		mockRepo := &mockProductRepositoryById{
			MockFindByIDProduct: nil,
			MockFindByIDError:   nil,
		}
		useCase := NewGetProductByIDUseCase(mockRepo)

		request := GetProductByIDRequest{ID: 999}

		foundProduct, err := useCase.Execute(ctx, request)

		if err != nil {
			t.Fatalf("expected no error, but got: %v", err)
		}
		if foundProduct != nil {
			t.Errorf("expected nil product, got %+v", foundProduct)
		}
	})

	t.Run("should return error when repository returns an error", func(t *testing.T) {
		expectedError := errors.New("database connection lost")

		mockRepo := &mockProductRepositoryById{
			MockFindByIDProduct: nil,
			MockFindByIDError:   expectedError,
		}
		useCase := NewGetProductByIDUseCase(mockRepo)

		request := GetProductByIDRequest{ID: 1}

		foundProduct, err := useCase.Execute(ctx, request)

		if err == nil {
			t.Fatalf("expected an error, but got none")
		}
		if !errors.Is(err, expectedError) {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
		if foundProduct != nil {
			t.Errorf("expected nil product on error, got %+v", foundProduct)
		}
	})
}
