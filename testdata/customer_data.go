package testdata

import (
	"github.com/VitoNaychev/bt-customer-svc/models/address_store"
	"github.com/VitoNaychev/bt-customer-svc/models/customer_store"
)

var PeterCustomer = customer_store.Customer{
	Id:          1,
	FirstName:   "Peter",
	LastName:    "Smith",
	PhoneNumber: "+359 88 576 5981",
	Email:       "petesmith@gmail.com",
	Password:    "firefirefire",
}

var AliceCustomer = customer_store.Customer{
	Id:          2,
	FirstName:   "Alice",
	LastName:    "Johnson",
	PhoneNumber: "+359 88 444 2222",
	Email:       "alicejohn@gmail.com",
	Password:    "helloJohn123",
}

var PeterAddress1 = address_store.Address{
	Id:           1,
	CustomerId:   1,
	Lat:          42.695111,
	Lon:          23.329184,
	AddressLine1: "Shipka Street 6",
	AddressLine2: "",
	City:         "Sofia",
	Country:      "Bulgaria",
}

var PeterAddress2 = address_store.Address{
	Id:           2,
	CustomerId:   1,
	Lat:          42.6938570,
	Lon:          23.3362452,
	AddressLine1: "ulitsa Gerogi S. Rakovski 96",
	AddressLine2: "",
	City:         "Sofia",
	Country:      "Bulgaria",
}

var AliceAddress = address_store.Address{
	Id:           3,
	CustomerId:   2,
	Lat:          42.6931204,
	Lon:          23.3225465,
	AddressLine1: "ut. Angel Kanchev 1",
	AddressLine2: "",
	City:         "Sofia",
	Country:      "Bulgaria",
}
