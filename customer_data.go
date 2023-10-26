package bt_customer_svc

var peterCustomer = Customer{
	Id:          0,
	FirstName:   "Peter",
	LastName:    "Smith",
	PhoneNumber: "+359 88 576 5981",
	Email:       "petesmith@gmail.com",
	Password:    "firefirefire",
}

var aliceCustomer = Customer{
	Id:          1,
	FirstName:   "Alice",
	LastName:    "Johnson",
	PhoneNumber: "+359 88 444 2222",
	Email:       "alicejohn@gmail.com",
	Password:    "helloJohn123",
}

var peterAddress1 = Address{
	Id:           0,
	CustomerId:   0,
	Lat:          42.695111,
	Lon:          23.329184,
	AddressLine1: "Shipka Street 6",
	AddressLine2: "",
	City:         "Sofia",
	Country:      "Bulgaria",
}

var peterAddress2 = Address{
	Id:           1,
	CustomerId:   0,
	Lat:          42.6938570,
	Lon:          23.3362452,
	AddressLine1: "ulitsa Gerogi S. Rakovski 96",
	AddressLine2: "",
	City:         "Sofia",
	Country:      "Bulgaria",
}

var aliceAddress = Address{
	Id:           2,
	CustomerId:   1,
	Lat:          42.6931204,
	Lon:          23.3225465,
	AddressLine1: "ut. Angel Kanchev 1",
	AddressLine2: "",
	City:         "Sofia",
	Country:      "Bulgaria",
}
