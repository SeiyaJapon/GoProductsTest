package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	HealthCheck(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	expected := "{\"status\":\"OK\"}\n"
	if rr.Body.String() != expected {
		t.Errorf("expected body %s, got %s", expected, rr.Body.String())
	}
}
