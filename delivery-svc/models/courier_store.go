package models

type CourierStore interface {
	CreateCourier(*Courier) error
	DeleteCourier(int) error
}
