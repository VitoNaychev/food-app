package testdata

import (
	"time"

	"github.com/VitoNaychev/food-app/delivery-svc/models"
)

var (
	VolenDelivery = models.Delivery{
		ID:                1,
		CourierID:         1,
		PickupAddressID:   1,
		DeliveryAddressID: 2,
		ReadyBy:           time.Time{},
		State:             models.CREATED,
	}

	PeterDelivery = models.Delivery{
		ID:                2,
		CourierID:         2,
		PickupAddressID:   3,
		DeliveryAddressID: 4,
		ReadyBy:           time.Time{},
		State:             models.IN_PROGRESS,
	}
)
