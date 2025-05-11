package router

import (
	"encoding/json"
	"net/http"
	"subito-cart/internal/errs"
	"subito-cart/internal/middleware"
	"subito-cart/internal/service"

	"github.com/gorilla/mux"
)

type orderRequest struct {
	Order struct {
		Items []service.OrderItem `json:"items" validate:"required"`
	} `json:"order" validate:"required"`
}

// Validate implements the RequestValidator interface
func (r orderRequest) Validate() error {
	return nil // Additional validation can be added here
}

type orderResponse struct {
	OrderID    int                  `json:"order_id"`
	OrderPrice float64              `json:"order_price"`
	OrderVAT   float64              `json:"order_vat"`
	Items      []service.PricedItem `json:"items"`
}

func RegisterOrderRoutes(r *mux.Router) {
	r.HandleFunc("/order", middleware.WithRequestValidation[orderRequest](handleOrder)).Methods("POST")
}

func handleOrder(w http.ResponseWriter, r *http.Request) {
	// Get the validated request from context
	req, ok := r.Context().Value(middleware.RequestKey).(orderRequest)
	if !ok {
		errs.SendErrorResponse(w, http.StatusInternalServerError, "Failed to process request", []errs.Error{
			{
				Field:   "request",
				Message: "Invalid request context",
				Code:    "CONTEXT_ERROR",
			},
		})
		return
	}

	//Business Logic
	orderID, totalPrice, totalVAT, pricedItems, err := service.CalculatePricing(req.Order.Items)
	if err != nil {
		errs.SendErrorResponse(w, http.StatusInternalServerError, "Failed to process order", []errs.Error{
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
