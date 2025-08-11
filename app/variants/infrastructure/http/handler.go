package http

import (
	"github.com/mytheresa/go-hiring-challenge/app/shared/api"
	"github.com/mytheresa/go-hiring-challenge/app/variants/domain"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	getProductByIDUseCase domain.GetProductByIDUseCaseInterface
}

func NewProductHandler(getProductByIDUseCase domain.GetProductByIDUseCaseInterface) *ProductHandler {
	return &ProductHandler{
		getProductByIDUseCase: getProductByIDUseCase,
	}
}

// GetProductById handles the HTTP GET request to retrieve a product by its ID along with its variants.
// It expects the ID to be passed as a query parameter.
func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	request := domain.GetProductByIDRequest{
		ID: uint(id),
	}

	product, err := h.getProductByIDUseCase.Execute(r.Context(), request)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, product)
}
