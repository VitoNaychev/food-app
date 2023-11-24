package models

type Address struct {
	Id           int
	CustomerId   int `db:"customer_id"`
	Lat          float64
	Lon          float64
	AddressLine1 string `db:"address_line1"`
	AddressLine2 string `db:"address_line2"`
	City         string
	Country      string
}
