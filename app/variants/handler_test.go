package variants

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockProductByIDUseCase struct {
	MockProduct *ProductWithVariants
	MockError   error
}

func (m *mockProductByIDUseCase) Execute(ctx context.Context, req GetProductByIDRequest) (*ProductWithVariants, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockProduct, nil
}

func TestGetProductById_Success(t *testing.T) {
	mockProduct := &ProductWithVariants{
		Product: &catalog.Product{
			ID:    1,
			Code:  "P001",
			Price: 50.0,
			Category: catalog.Category{
				ID:   "1",
				Code: "Accessories",
				Name: "Accesorios",
			},
		},
		Variants: []Variant{
			{
				ID:    "1",
				Name:  "Variant 1",
				SKU:   "SKU001",
				Price: 50.0,
			},
		},
	}

	mockUseCase := &mockProductByIDUseCase{
		MockProduct: mockProduct,
	}
	handler := NewProductHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/catalog?id=1", nil)
	rr := httptest.NewRecorder()

	handler.GetProductById(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var receivedProduct ProductWithVariants
	if err := json.Unmarshal(rr.Body.Bytes(), &receivedProduct); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if receivedProduct.Product.ID != 1 {
		t.Errorf("expected product_variants ID 1, got %d", receivedProduct.Product.ID)
	}
}

func TestGetProductById_InvalidID(t *testing.T) {
	handler := NewProductHandler(&mockProductByIDUseCase{})

	req := httptest.NewRequest(http.MethodGet, "/catalog/abc", nil)
	rr := httptest.NewRecorder()

	handler.GetProductById(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetProductById_ErrorFromUseCase(t *testing.T) {
	mockUseCase := &mockProductByIDUseCase{
		MockError: errors.New("product_variants not found in use case"),
	}
	handler := NewProductHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/catalog?id=1", nil)
	rr := httptest.NewRecorder()

	handler.GetProductById(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
