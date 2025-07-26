package handlers_variants

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/domain/product"
	"github.com/mytheresa/go-hiring-challenge/app/variants/application/usecases_variants"
	"github.com/mytheresa/go-hiring-challenge/app/variants/domain/product_variants"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type mockProductByIDUseCase struct {
	MockProduct *product.Product
	MockError   error
}

func (m *mockProductByIDUseCase) Execute(ctx context.Context, req usecases_variants.GetProductByIDRequest) (*product.Product, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockProduct, nil
}

func TestGetProductById_Success(t *testing.T) {
	mockProduct := &product.Product{
		ID:    1,
		Code:  "P001",
		Price: 50.0,
		Category: product.Category{
			ID:   strconv.Itoa(1),
			Code: "Accessories",
			Name: "Accesorios",
		},
		Variants: []product_variants.Variant{},
	}

	mockUseCase := &mockProductByIDUseCase{
		MockProduct: mockProduct,
		MockError:   nil,
	}
	handler := NewProductHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/catalog/1", nil)
	rr := httptest.NewRecorder()

	handler.GetProductById(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var receivedProduct product.Product
	if err := json.Unmarshal(rr.Body.Bytes(), &receivedProduct); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if receivedProduct.ID != 1 {
		t.Errorf("expected product_variants ID 1, got %d", receivedProduct.ID)
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
		MockProduct: nil,
		MockError:   errors.New("product_variants not found in use case"),
	}
	handler := NewProductHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/catalog/1", nil)
	rr := httptest.NewRecorder()

	handler.GetProductById(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
