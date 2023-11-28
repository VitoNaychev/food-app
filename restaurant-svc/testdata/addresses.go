package testdata

import (
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

var ShackAddress = models.Address{
	ID:           1,
	RestaurantID: 1,
	Lat:          42.6359749959353,
	Lon:          23.3807774069591,
	AddressLine1: "ul. Filip Avramov 411, gk Mladost 4",
	AddressLine2: "",
	City:         "Sofia",
	Country:      "Bulgaria",
}

var DominosAddress = models.Address{
	ID:           2,
	RestaurantID: 2,
	Lat:          42.6362464985259,
	Lon:          23.3698686256139,
	AddressLine1: "Aleksandar Malinov Boulevard 78",
	AddressLine2: "",
	City:         "Sofia",
	Country:      "Bulgaria",
}
