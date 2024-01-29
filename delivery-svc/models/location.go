package models

type Location struct {
	CourierID int `db:"courier_id"`
	Lat       float32
	Lon       float32
}
