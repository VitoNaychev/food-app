package testdata

import "github.com/VitoNaychev/bt-order-svc/models"

var (
	ChickenShackAddress = models.Address{
		ID:           1,
		Lat:          42.635934305,
		Lon:          23.380761684,
		AddressLine1: "ul. Filip Avramov 411",
		AddressLine2: "",
		City:         "Sofia",
		Country:      "Bulgaria",
	}
	PeterAddress1 = models.Address{
		ID:           2,
		Lat:          42.695111,
		Lon:          23.329184,
		AddressLine1: "Shipka Street 6",
		AddressLine2: "",
		City:         "Sofia",
		Country:      "Bulgaria",
	}
	PeterAddress2 = models.Address{
		ID:           3,
		Lat:          42.6938570,
		Lon:          23.3362452,
		AddressLine1: "ulitsa Gerogi S. Rakovski 96",
		AddressLine2: "",
		City:         "Sofia",
		Country:      "Bulgaria",
	}
	AliceAddress = models.Address{
		ID:           4,
		Lat:          42.6931204,
		Lon:          23.3225465,
		AddressLine1: "ut. Angel Kanchev 1",
		AddressLine2: "",
		City:         "Sofia",
		Country:      "Bulgaria",
	}
)
