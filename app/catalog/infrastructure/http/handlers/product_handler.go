package handlers

import (
	"github.com/mytheresa/go-hiring-challenge/app/catalog/application/usecases"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/infrastructure/http/api"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	getProductByIDUseCase usecases.GetProductByIDUseCaseInterface
}

func NewProductHandler(getProductByIDUseCase usecases.GetProductByIDUseCaseInterface) *ProductHandler {
	return &ProductHandler{
		getProductByIDUseCase: getProductByIDUseCase,
	}
}

func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	request := usecases.GetProductByIDRequest{
		ID: uint(id),
	}

	product, err := h.getProductByIDUseCase.Execute(r.Context(), request)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, product)
}
