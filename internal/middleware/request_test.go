package middleware

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockRequest implements RequestValidator
type mockRequest struct {
	Name string `json:"name" validate:"required"`
}

func (m mockRequest) Validate() error {
	return nil
}

// testHandler captures if it's called and reads from context
func testHandlerCalled(t *testing.T, expectName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val, ok := r.Context().Value(requestKey).(mockRequest)
		if !ok {
			t.Errorf("expected mockRequest in context, got none")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if val.Name != expectName {
			t.Errorf("expected name %q, got %q", expectName, val.Name)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func TestWithRequestValidation_ValidInput(t *testing.T) {
	body := `{"name":"Alice"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler := WithRequestValidation[mockRequest](testHandlerCalled(t, "Alice"))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}
}

func TestWithRequestValidation_InvalidJSON(t *testing.T) {
	body := `{ invalid json }`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler := WithRequestValidation[mockRequest](http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called on invalid JSON")
	}))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestWithRequestValidation_MissingRequiredField(t *testing.T) {
	body := `{}` // missing "name"
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler := WithRequestValidation[mockRequest](http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called when validation fails")
	}))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestWithRequestValidation_BodyReadError(t *testing.T) {
	brokenReader := io.NopCloser(&errorReader{})
	req := httptest.NewRequest(http.MethodPost, "/", brokenReader)
	rec := httptest.NewRecorder()

	handler := WithRequestValidation[mockRequest](http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called on read error")
	}))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rec.Code)
	}
}

// errorReader simulates a body read failure
type errorReader struct{}

func (e *errorReader) Read([]byte) (int, error) {
	return 0, errors.New("simulated read error")
}
