package http

import (
	"github.com/mytheresa/go-hiring-challenge/app/shared/api"
	"net/http"
)

type Handler struct {
	Status string `json:"status"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	response := Handler{Status: "OK"}

	api.OKResponse(w, response)
}
