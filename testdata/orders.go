package testdata

import (
	"time"

	"github.com/VitoNaychev/bt-order-svc/models"
)

var (
	PeterOrder = models.Order{
		ID:           1,
		CustomerID:   1,
		RestaurantID: 1,
		Items:        []int{1, 2, 3},
		Total:        13.12,
		DeliveryTime: time.Date(2024, 2, 12, 18, 30, 00, 00, time.Local),
	}
	AliceOrder = models.Order{
		ID:           2,
		CustomerID:   2,
		RestaurantID: 1,
		Items:        []int{2, 2, 5},
		Total:        14.42,
		DeliveryTime: time.Date(2024, 2, 12, 17, 00, 00, 00, time.Local),
	}
)
