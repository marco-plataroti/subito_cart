package error

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Errors  []Error `json:"errors,omitempty"`
}

// Error represents a single validation error
type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// sendErrorResponse sends a standardized error response
func SendErrorResponse(w http.ResponseWriter, status int, message string, errors []Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  status,
		Message: message,
		Errors:  errors,
	})
}