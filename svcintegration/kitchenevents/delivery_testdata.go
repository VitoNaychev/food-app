package kitchenevents

import (
	"time"

	"github.com/VitoNaychev/food-app/delivery-svc/models"
)

var (
	volenPickupAddress = models.Address{
		ID:           1,
		Lat:          42.6359749959353,
		Lon:          23.3807774069591,
		AddressLine1: "ul. Filip Avramov 411, gk Mladost 4",
		AddressLine2: "",
		City:         "Sofia",
		Country:      "Bulgaria",
	}

	volenDeliveryAddress = models.Address{
		ID:           2,
		Lat:          42.695111,
		Lon:          23.329184,
		AddressLine1: "Shipka Street 6",
		AddressLine2: "",
		City:         "Sofia",
		Country:      "Bulgaria",
	}

	volenCourier = models.Courier{
		ID:   1,
		Name: "Volen",
	}

	volenDelivery = models.Delivery{
		ID:                1,
		CourierID:         1,
		PickupAddressID:   1,
		DeliveryAddressID: 2,
		ReadyBy:           time.Time{},
		State:             models.PENDING,
	}
)
