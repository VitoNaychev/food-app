package models

type CourierStore interface {
	DeleteCourier(id int) error
	UpdateCourier(*Courier) error
	CreateCourier(*Courier) error
	GetCourierByID(id int) (Courier, error)
	GetCourierByEmail(email string) (Courier, error)
}
