package service

import (
	"math/rand"
	"time"

	"subito-cart/internal/repository"
)

type OrderItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type PricedItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	VAT       float64 `json:"vat"`
}

func CalculatePricing(items []OrderItem) (int, float64, float64, []PricedItem, error) {
	var (
		totalPrice float64
		totalVAT   float64
		result     []PricedItem
	)

	for _, item := range items {
		product, err := repository.GetProductByID(item.ProductID)
		if err != nil {
			return 0, 0, 0, nil, err
		}

		itemPrice := float64(item.Quantity) * product.Price
		itemVAT := itemPrice * product.VAT

		totalPrice += itemPrice
		totalVAT += itemVAT

		result = append(result, PricedItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     itemPrice,
			VAT:       itemVAT,
		})
	}

	orderID := generateOrderID()
	return orderID, round(totalPrice), round(totalVAT), result, nil
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateOrderID() int {
	return rng.Intn(1_000_000_000)
}

func round(f float64) float64 {
	return float64(int(f*100+0.5)) / 100
}
