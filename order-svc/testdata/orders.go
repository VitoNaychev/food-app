package testdata

import (
	"time"

	"github.com/VitoNaychev/food-app/order-svc/models"
)

var (
	PeterCreatedOrder = models.Order{
		ID:              1,
		CustomerID:      1,
		RestaurantID:    1,
		Total:           13.12,
		DeliveryTime:    time.Date(2024, 2, 12, 18, 30, 00, 00, time.UTC),
		PickupAddress:   1,
		DeliveryAddress: 2,
		Status:          models.APPROVAL_PENDING,
	}
	PeterCompletedOrder = models.Order{
		ID:              2,
		CustomerID:      1,
		RestaurantID:    1,
		Total:           19.22,
		DeliveryTime:    time.Date(2022, 2, 12, 18, 30, 00, 00, time.UTC),
		PickupAddress:   1,
		DeliveryAddress: 3,
		Status:          models.COMPLETED,
	}
	AliceOrder = models.Order{
		ID:              3,
		CustomerID:      2,
		RestaurantID:    1,
		Total:           14.42,
		DeliveryTime:    time.Date(2022, 2, 12, 17, 00, 00, 00, time.UTC),
		PickupAddress:   1,
		DeliveryAddress: 4,
		Status:          models.COMPLETED,
	}
)
