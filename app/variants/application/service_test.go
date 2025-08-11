package application

import (
	"context"
	"errors"
	domaincatalog "github.com/mytheresa/go-hiring-challenge/app/catalog/domain"
	"github.com/mytheresa/go-hiring-challenge/app/variants/domain"
	"strconv"
	"testing"
)

type mockVariantRepository struct {
	MockFindByIDProduct *domain.ProductWithVariants
	MockFindByIDError   error

	MockFindVariantsByProductID []domain.Variant
}

func (m *mockVariantRepository) FindByID(ctx context.Context, id uint) (*domain.ProductWithVariants, error) {
	if m.MockFindByIDError != nil {
		return nil, m.MockFindByIDError
	}

	if m.MockFindByIDProduct == nil {
		return nil, nil
	}

	return &domain.ProductWithVariants{
		Product:  m.MockFindByIDProduct.Product,
		Variants: m.MockFindByIDProduct.Variants,
	}, nil
}

func (m *mockVariantRepository) FindVariantsByProductID(ctx context.Context, productID uint) ([]domain.Variant, error) {
	return m.MockFindVariantsByProductID, nil
}

func TestGetProductByIDUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should return product_variants when repository finds it", func(t *testing.T) {
		expectedProduct := &domain.ProductWithVariants{
			Product: &domaincatalog.Product{
				ID:    1,
				Code:  "P001",
				Price: 75.0,
				Category: domaincatalog.Category{
					ID:   strconv.Itoa(201),
					Code: "SHOES",
					Name: "Shoes",
				},
			},
			Variants: []domain.Variant{},
		}

		mockVariantRepo := &mockVariantRepository{
			MockFindByIDProduct: &domain.ProductWithVariants{
				Product:  expectedProduct.Product,
				Variants: expectedProduct.Variants,
			},
			MockFindByIDError:           nil,
			MockFindVariantsByProductID: expectedProduct.Variants,
		}

		useCase := NewGetProductByIDUseCase(mockVariantRepo)

		request := domain.GetProductByIDRequest{ID: 1}

		foundProduct, err := useCase.Execute(ctx, request)

		if err != nil {
			t.Fatalf("expected no error, but got: %v", err)
		}
		if foundProduct == nil {
			t.Fatal("expected a product_variants, got nil")
		}
		if foundProduct.Product.ID != expectedProduct.Product.ID ||
			foundProduct.Product.Code != expectedProduct.Product.Code {
			t.Errorf("expected product_variants %+v, got %+v", expectedProduct, foundProduct)
		}
	})

	t.Run("should return nil and no error when product_variants is not found", func(t *testing.T) {
		mockVariantRepo := &mockVariantRepository{
			MockFindByIDProduct:         nil,
			MockFindByIDError:           nil,
			MockFindVariantsByProductID: []domain.Variant{},
		}
		useCase := NewGetProductByIDUseCase(mockVariantRepo)

		request := domain.GetProductByIDRequest{ID: 999}

		foundProduct, err := useCase.Execute(ctx, request)

		if err != nil {
			t.Fatalf("expected no error, but got: %v", err)
		}
		if foundProduct != nil {
			t.Errorf("expected nil product_variants, got %+v", foundProduct)
		}
	})

	t.Run("should return error when repository returns an error", func(t *testing.T) {
		expectedError := errors.New("database connection lost")

		mockVariantRepo := &mockVariantRepository{
			MockFindByIDProduct:         nil,
			MockFindByIDError:           expectedError,
			MockFindVariantsByProductID: []domain.Variant{},
		}
		useCase := NewGetProductByIDUseCase(mockVariantRepo)

		request := domain.GetProductByIDRequest{ID: 1}

		foundProduct, err := useCase.Execute(ctx, request)

		if err == nil {
			t.Fatalf("expected an error, but got none")
		}
		if !errors.Is(err, expectedError) {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
		if foundProduct != nil {
			t.Errorf("expected nil product_variants on error, got %+v", foundProduct)
		}
	})
}
