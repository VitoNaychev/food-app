package address_store

type Address struct {
	Id           int
	CustomerId   int
	Lat          float64
	Lon          float64
	AddressLine1 string
	AddressLine2 string
	City         string
	Country      string
}
