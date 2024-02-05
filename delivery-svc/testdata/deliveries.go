package testdata

import (
	"github.com/VitoNaychev/food-app/delivery-svc/models"
)

var (
	VolenDelivery = models.Delivery{
		ID:                1,
		CourierID:         1,
		PickupAddressID:   1,
		DeliveryAddressID: 2,
		ReadyBy:           models.ZeroTime,
		State:             models.PENDING,
	}

	VolenActiveDelivery = models.Delivery{
		ID:                1,
		CourierID:         1,
		PickupAddressID:   1,
		DeliveryAddressID: 2,
		ReadyBy:           models.ZeroTime,
		State:             models.READY_FOR_PICKUP,
	}

	PeterDelivery = models.Delivery{
		ID:                2,
		CourierID:         2,
		PickupAddressID:   3,
		DeliveryAddressID: 4,
		ReadyBy:           models.ZeroTime,
		State:             models.IN_PROGRESS,
	}

	AliceDelivery = models.Delivery{
		ID:                3,
		CourierID:         3,
		PickupAddressID:   5,
		DeliveryAddressID: 6,
		ReadyBy:           models.ZeroTime,
		State:             models.READY_FOR_PICKUP,
	}

	JohnDelivery = models.Delivery{
		ID:                4,
		CourierID:         4,
		PickupAddressID:   7,
		DeliveryAddressID: 8,
		ReadyBy:           models.ZeroTime,
		State:             models.ON_ROUTE,
	}

	IvoDelivery = models.Delivery{
		ID:                5,
		CourierID:         5,
		PickupAddressID:   9,
		DeliveryAddressID: 10,
		ReadyBy:           models.ZeroTime,
		State:             models.COMPLETED,
	}
)
