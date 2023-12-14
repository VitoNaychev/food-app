package testdata

import "github.com/VitoNaychev/food-app/courier-svc/models"

var (
	MichaelCourier = models.Courier{
		ID:          1,
		FirstName:   "Michael",
		LastName:    "Scott",
		PhoneNumber: "+359 88 999 1111",
		Email:       "mscarn@gmail.com",
		Password:    "samplepassword",
		IBAN:        "DE89370400440532013000",
	}

	JimCourier = models.Courier{
		ID:          2,
		FirstName:   "Jim",
		LastName:    "Halpert",
		PhoneNumber: "+359 88 222 3333",
		Email:       "jimhalp@dominos.com",
		Password:    "samplepassword",
		IBAN:        "DE89370400440532013000",
	}
)
