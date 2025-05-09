package router

import (
	"encoding/json"
	"net/http"
	"subito-cart/internal/service"

	"github.com/gorilla/mux"
)

type orderRequest struct {
	Order struct {
		Items []service.OrderItem `json:"items"`
	} `json:"order"`
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
	var req orderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	orderID, totalPrice, totalVAT, pricedItems, err := service.CalculatePricing(req.Order.Items)
	if err != nil {
		http.Error(w, "Failed to calculate order", http.StatusInternalServerError)
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
