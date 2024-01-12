package testdata

import "github.com/VitoNaychev/food-app/delivery-svc/models"

var (
	VolenPickupAddress = models.Address{
		ID:           1,
		Lat:          42.6359749959353,
		Lon:          23.3807774069591,
		AddressLine1: "ul. Filip Avramov 411, gk Mladost 4",
		AddressLine2: "",
		City:         "Sofia",
		Country:      "Bulgaria",
	}

	VolenDeliveryAddress = models.Address{
		ID:           2,
		Lat:          42.695111,
		Lon:          23.329184,
		AddressLine1: "Shipka Street 6",
		AddressLine2: "",
		City:         "Sofia",
		Country:      "Bulgaria",
	}
)
