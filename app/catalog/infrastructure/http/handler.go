package http

import (
	"github.com/mytheresa/go-hiring-challenge/app/catalog/domain"
	"github.com/mytheresa/go-hiring-challenge/app/shared/api"
	"net/http"
	"strconv"
)

type CatalogHandler struct {
	getCatalogUseCase domain.GetCatalogUseCaseInterface
}

func NewCatalogHandler(getCatalogUseCase domain.GetCatalogUseCaseInterface) *CatalogHandler {
	return &CatalogHandler{
		getCatalogUseCase: getCatalogUseCase,
	}
}

// GetCatalog handles the HTTP GET request to retrieve the product catalog.
func (h *CatalogHandler) GetCatalog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	category := query.Get("category")
	priceLtStr := query.Get("price_lt")

	var priceLt *float64
	if priceLtStr != "" {
		price, err := strconv.ParseFloat(priceLtStr, 64)
		if err == nil {
			priceLt = &price
		}
	}

	request := domain.GetCatalogRequest{
		Offset:   offset,
		Limit:    limit,
		Category: category,
		PriceLt:  priceLt,
	}

	products, err := h.getCatalogUseCase.Execute(r.Context(), request)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, products)
}
