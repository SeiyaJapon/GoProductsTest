package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/application/usecases"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/domain/product"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type mockCatalogUseCase struct {
	MockProducts []product.Product
	MockError    error
}

func (m *mockCatalogUseCase) Execute(ctx context.Context, req usecases.GetCatalogRequest) ([]product.Product, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockProducts, nil
}

func TestGetCatalog_Success(t *testing.T) {
	mockUseCase := &mockCatalogUseCase{
		MockProducts: []product.Product{
			{
				ID:    1,
				Code:  "P001",
				Price: 150.0,
				Category: product.Category{
					ID:   strconv.Itoa(101),
					Code: "Clothing",
					Name: "Ropa",
				},
				Variants: []product.Variant{},
			},
		},
		MockError: nil,
	}
	handler := NewCatalogHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/catalog?offset=0&limit=10&category=Clothing&price_lt=200", nil)
	rr := httptest.NewRecorder()

	handler.GetCatalog(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var products []product.Product
	if err := json.Unmarshal(rr.Body.Bytes(), &products); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(products) != 1 || products[0].Category.Code != "Clothing" {
		t.Errorf("unexpected response: %+v", products)
	}
}

func TestGetCatalog_InternalError(t *testing.T) {
	mockUseCase := &mockCatalogUseCase{
		MockProducts: nil,
		MockError:    errors.New("something went wrong in use case"),
	}
	handler := NewCatalogHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
	rr := httptest.NewRecorder()

	handler.GetCatalog(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
