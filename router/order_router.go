package router

import (
	"encoding/json"
	"io"
	"net/http"
	"subito-cart/internal/error"
	"subito-cart/internal/service"
	"subito-cart/internal/validator"

	"github.com/gorilla/mux"
)

type orderRequest struct {
	Order struct {
		Items []service.OrderItem `json:"items" validate:"required"`
	} `json:"order" validate:"required"`
}

type orderResponse struct {
	OrderID    int                  `json:"order_id"`
	OrderPrice float64              `json:"order_price"`
	OrderVAT   float64              `json:"order_vat"`
	Items      []service.PricedItem `json:"items"`
}

func RegisterOrderRoutes(r *mux.Router) {
	r.HandleFunc("/order", handleOrder).Methods("POST")
}

func handleOrder(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		error.SendErrorResponse(w, http.StatusBadRequest, "Failed to read request body", []error.Error{
			{
				Field:   "body",
				Message: "Could not read request body",
				Code:    "READ_ERROR",
			},
		})
		return
	}

	var req orderRequest
	// Validate JSON and struct
	if err := validator.ValidateJSON(body, &req); err != nil {
		var validationErrors []error.Error
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				validationErrors = append(validationErrors, error.Error{
					Field:   e.Field,
					Message: e.Message,
					Code:    e.Tag,
				})
			}
		}
		error.SendErrorResponse(w, http.StatusBadRequest, "Invalid request", validationErrors)
		return
	}

	//Business Logic
	orderID, totalPrice, totalVAT, pricedItems, err := service.CalculatePricing(req.Order.Items)
	if err != nil {
		error.SendErrorResponse(w, http.StatusInternalServerError, "Failed to process order", []error.Error{
			{
				Field:   "order",
				Message: err.Error(),
				Code:    "PROCESSING_ERROR",
			},
		})
		return
	}

	resp := orderResponse{
		OrderID:    orderID,
		OrderPrice: totalPrice,
		OrderVAT:   totalVAT,
		Items:      pricedItems,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
