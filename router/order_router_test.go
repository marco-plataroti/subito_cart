package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// mock data
var validRequestBody = []byte(`{
	"order": {
		"items": [
			{ "product_id": 1, "quantity": 2 },
			{ "product_id": 2, "quantity": 1 }
		]
	}
}`)

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	RegisterOrderRoutes(r)
	return r
}

func TestHandleOrder_MissingFields(t *testing.T) {
	invalidBody := []byte(`{"order": {}}`) // items field is missing

	router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(invalidBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestHandleOrder_InvalidJSON(t *testing.T) {
	badJSON := []byte(`{ this is invalid json }`)

	router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(badJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestHandleOrder_MissingContextInjection(t *testing.T) {
	// Manually bypass middleware to simulate missing context
	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(validRequestBody))
	req = req.WithContext(req.Context()) // no injected request

	rec := httptest.NewRecorder()
	handleOrder(rec, req) // directly call the handler

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 Internal Server Error, got %d", rec.Code)
	}
}
