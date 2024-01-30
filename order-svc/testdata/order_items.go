package testdata

import "github.com/VitoNaychev/food-app/order-svc/models"

var (
	PeterCreatedOrderItems = []models.OrderItem{
		{
			ID:         1,
			OrderID:    1,
			MenuItemID: 1,
			Quantity:   1,
		},
		{
			ID:         2,
			OrderID:    1,
			MenuItemID: 2,
			Quantity:   1,
		},
		{
			ID:         3,
			OrderID:    1,
			MenuItemID: 3,
			Quantity:   1,
		},
	}

	PeterCompletedOrderItems = []models.OrderItem{
		{
			ID:         4,
			OrderID:    2,
			MenuItemID: 3,
			Quantity:   4,
		},
		{
			ID:         5,
			OrderID:    2,
			MenuItemID: 5,
			Quantity:   1,
		},
	}

	AliceOrderItems = []models.OrderItem{
		{
			ID:         6,
			OrderID:    3,
			MenuItemID: 2,
			Quantity:   2,
		},
		{
			ID:         7,
			OrderID:    3,
			MenuItemID: 5,
			Quantity:   1,
		},
	}
)
