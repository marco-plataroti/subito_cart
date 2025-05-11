package middleware

import (
	"context"
	"io"
	"net/http"
	"subito-cart/internal/errs"
	"subito-cart/internal/validator"
)

// RequestValidator is a generic type constraint for request structs
type RequestValidator interface {
	Validate() error
}

type contextKey string

const RequestKey = contextKey("validated_request")

// WithRequestValidation creates a middleware that validates the request body
func WithRequestValidation[T RequestValidator](next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			errs.SendErrorResponse(w, http.StatusBadRequest, "Failed to read request body", []errs.Error{
				{
					Field:   "body",
					Message: "Could not read request body",
					Code:    "READ_ERROR",
				},
			})
			return
		}

		// Create a new instance of the request type
		var req T

		// Validate JSON and struct
		if err := validator.ValidateJSON(body, &req); err != nil {
			var validationErrors []errs.Error
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, e := range ve {
					validationErrors = append(validationErrors, errs.Error{
						Field:   e.Field,
						Message: e.Message,
						Code:    e.Tag,
					})
				}
			}
			errs.SendErrorResponse(w, http.StatusBadRequest, "Invalid request", validationErrors)
			return
		}

		// Store the validated request in the context
		ctx := context.WithValue(r.Context(), RequestKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
