package testdata

import "github.com/VitoNaychev/food-app/restaurant-svc/models"

var (
	Restaurant = models.Restaurant{
		ID:          1,
		Name:        "Chicken Shack",
		PhoneNumber: "+359 567 0890",
		Email:       "shack@gmail.com",
		Password:    "samplepassword",
		IBAN:        "DE89370400440532013000",
		Status:      models.CREATION_PENDING,
	}
)
