package repository

import "errors"

type Product struct {
	ID    int
	Price float64 // Price per unit
	VAT   float64 // VAT rate as a fraction, e.g., 0.10 for 10%
}

var products = map[int]Product{
	1: {ID: 1, Price: 2.00, VAT: 0.10}, // 10% VAT
	2: {ID: 2, Price: 1.50, VAT: 0.10}, // 10% VAT
	3: {ID: 3, Price: 3.00, VAT: 0.10}, // 10% VAT
}

func GetProductByID(id int) (*Product, error) {
	if product, ok := products[id]; ok {
		return &product, nil
	}
	return nil, errors.New("product not found")
}

