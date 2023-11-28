package testdata

import "github.com/VitoNaychev/food-app/restaurant-svc/models"

var (
	ShackRestaurant = models.Restaurant{
		ID:          1,
		Name:        "Chicken Shack",
		PhoneNumber: "+359 567 0890",
		Email:       "shack@gmail.com",
		Password:    "samplepassword",
		IBAN:        "DE89370400440532013000",
		Status:      models.CREATION_PENDING,
	}

	DominosRestaurant = models.Restaurant{
		ID:          2,
		Name:        "Dominos",
		PhoneNumber: "+359 88 553 1234",
		Email:       "pizza@dominos.com",
		Password:    "samplepassword",
		IBAN:        "DE89370400440532013000",
		Status:      models.VALID,
	}
)
