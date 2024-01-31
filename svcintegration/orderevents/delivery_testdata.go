package orderevents

import (
	"time"

	"github.com/VitoNaychev/food-app/delivery-svc/models"
)

var (
	volenCourier = models.Courier{
		ID:   1,
		Name: "Volen",
	}

	volenPickupAddress = models.Address{
		ID:           1,
		Lat:          42.635934305,
		Lon:          23.380761684,
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

	volenDelivery = models.Delivery{
		ID:                1,
		CourierID:         1,
		PickupAddressID:   1,
		DeliveryAddressID: 2,
		ReadyBy:           time.Time{},
		State:             models.PENDING,
	}
)
