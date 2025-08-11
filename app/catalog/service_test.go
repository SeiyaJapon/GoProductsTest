package catalog

import (
	"context"
	"errors"
	"strconv"
	"testing"
)

type mockProductRepository struct {
	MockFindAllProducts []Product
	MockFindAllError    error
	MockFindByIDProduct *Product
}

func (m *mockProductRepository) FindAll(ctx context.Context, offset int, limit int, category string, priceLt *float64) ([]Product, error) {
	return m.MockFindAllProducts, m.MockFindAllError
}

func (m *mockProductRepository) FindByID(ctx context.Context, id uint) (*Product, error) {
	return nil, errors.New("FindByID not implemented in mock for this test")
}

func TestGetCatalogUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	floatPtr := func(f float64) *float64 { return &f }

	t.Run("should return products when repository finds them", func(t *testing.T) {
		expectedProducts := []Product{
			{
				ID:    1,
				Code:  "P001",
				Price: 100.0,
				Category: Category{
					ID:   strconv.Itoa(101),
					Code: "ELECTRONICS",
					Name: "Electronics",
				},
			},
		}

		mockRepo := &mockProductRepository{
			MockFindAllProducts: expectedProducts,
			MockFindAllError:    nil,
		}

		useCase := NewGetCatalogUseCase(mockRepo)

		request := GetCatalogRequest{
			Offset:   0,
			Limit:    10,
			Category: "Electronics",
			PriceLt:  floatPtr(150.0),
		}

		products, err := useCase.Execute(ctx, request)

		if err != nil {
			t.Fatalf("expected no error, but got: %v", err)
		}
		if len(products) != len(expectedProducts) {
			t.Fatalf("expected %d products, got %d", len(expectedProducts), len(products))
		}
		if products[0].ID != expectedProducts[0].ID || products[0].Code != expectedProducts[0].Code {
			t.Errorf("expected product_variants %+v, got %+v", expectedProducts[0], products[0])
		}
	})

	t.Run("should return empty slice when no products are found", func(t *testing.T) {
		mockRepo := &mockProductRepository{
			MockFindAllProducts: []Product{},
			MockFindAllError:    nil,
		}
		useCase := NewGetCatalogUseCase(mockRepo)

		request := GetCatalogRequest{Offset: 0, Limit: 10, Category: "Books", PriceLt: nil}

		products, err := useCase.Execute(ctx, request)

		if err != nil {
			t.Fatalf("expected no error, but got: %v", err)
		}
		if len(products) != 0 {
			t.Errorf("expected empty slice, got %d products", len(products))
		}
	})

	t.Run("should return error when repository returns an error", func(t *testing.T) {
		expectedError := errors.New("database connection failed")

		mockRepo := &mockProductRepository{
			MockFindAllProducts: nil,
			MockFindAllError:    expectedError,
		}
		useCase := NewGetCatalogUseCase(mockRepo)

		request := GetCatalogRequest{Offset: 0, Limit: 10, Category: "Tools", PriceLt: floatPtr(50.0)}

		products, err := useCase.Execute(ctx, request)

		if err == nil {
			t.Fatalf("expected an error, but got none")
		}
		if !errors.Is(err, expectedError) {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
		if products != nil {
			t.Errorf("expected nil products on error, got %+v", products)
		}
	})
}
