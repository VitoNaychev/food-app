package testdata

import (
	"time"

	"github.com/VitoNaychev/bt-order-svc/models"
)

var (
	PeterOrder1 = models.Order{
		ID:              1,
		CustomerID:      1,
		RestaurantID:    1,
		Items:           []int{1, 2, 3},
		Total:           13.12,
		DeliveryTime:    time.Date(2024, 2, 12, 18, 30, 00, 00, time.UTC),
		PickupAddress:   1,
		DeliveryAddress: 2,
		Status:          models.APPROVAL_PENDING,
	}
	PeterOrder2 = models.Order{
		ID:              1,
		CustomerID:      1,
		RestaurantID:    1,
		Items:           []int{3, 3, 3, 3, 5},
		Total:           19.22,
		DeliveryTime:    time.Date(2022, 2, 12, 18, 30, 00, 00, time.UTC),
		PickupAddress:   1,
		DeliveryAddress: 3,
		Status:          models.COMPLETED,
	}
	AliceOrder = models.Order{
		ID:              2,
		CustomerID:      2,
		RestaurantID:    1,
		Items:           []int{2, 2, 5},
		Total:           14.42,
		DeliveryTime:    time.Date(2022, 2, 12, 17, 00, 00, 00, time.UTC),
		PickupAddress:   1,
		DeliveryAddress: 4,
		Status:          models.COMPLETED,
	}
)
