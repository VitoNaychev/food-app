package testdata

import (
	"github.com/VitoNaychev/food-app/order-svc/models"
)

var (
	PeterCreatedOrder = models.Order{
		ID:              1,
		CustomerID:      1,
		RestaurantID:    1,
		Total:           13.12,
		PickupAddress:   1,
		DeliveryAddress: 2,
		Status:          models.APPROVAL_PENDING,
	}
	PeterCompletedOrder = models.Order{
		ID:              2,
		CustomerID:      1,
		RestaurantID:    1,
		Total:           19.22,
		PickupAddress:   1,
		DeliveryAddress: 3,
		Status:          models.COMPLETED,
	}
	AliceOrder = models.Order{
		ID:              3,
		CustomerID:      2,
		RestaurantID:    1,
		Total:           14.42,
		PickupAddress:   1,
		DeliveryAddress: 4,
		Status:          models.COMPLETED,
	}
)
