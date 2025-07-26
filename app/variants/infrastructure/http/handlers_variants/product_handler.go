package handlers_variants

import (
	"github.com/mytheresa/go-hiring-challenge/app/shared/infrastructure/http/api"
	"github.com/mytheresa/go-hiring-challenge/app/variants/application/usecases_variants"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	getProductByIDUseCase usecases_variants.GetProductByIDUseCaseInterface
}

func NewProductHandler(getProductByIDUseCase usecases_variants.GetProductByIDUseCaseInterface) *ProductHandler {
	return &ProductHandler{
		getProductByIDUseCase: getProductByIDUseCase,
	}
}

func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 3 || parts[1] != "catalog" {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid URL format. Expected /catalog/{id}")
		return
	}

	idStr := parts[2]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	request := usecases_variants.GetProductByIDRequest{
		ID: uint(id),
	}

	product, err := h.getProductByIDUseCase.Execute(r.Context(), request)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, product)
}
