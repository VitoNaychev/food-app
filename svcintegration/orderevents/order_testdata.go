package orderevents

import (
	"github.com/VitoNaychev/food-app/order-svc/models"
)

var (
	peterCustomerID = 1

	peterCreatedOrder = models.Order{
		ID:              1,
		CustomerID:      1,
		RestaurantID:    1,
		Total:           13.12,
		PickupAddress:   1,
		DeliveryAddress: 2,
		Status:          models.APPROVAL_PENDING,
	}

	peterCreatedOrderItems = []models.OrderItem{
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

	peterAddress1 = models.Address{
		ID:           2,
		Lat:          42.695111,
		Lon:          23.329184,
		AddressLine1: "Shipka Street 6",
		AddressLine2: "",
		City:         "Sofia",
		Country:      "Bulgaria",
	}

	chickenShackAddress = models.Address{
		ID:           1,
		Lat:          42.635934305,
		Lon:          23.380761684,
		AddressLine1: "ul. Filip Avramov 411, gk Mladost 4",
		AddressLine2: "",
		City:         "Sofia",
		Country:      "Bulgaria",
	}
)
