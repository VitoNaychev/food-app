package models

type DeliveryStore interface {
	GetDeliveryByID(int) (Delivery, error)
	UpdateDelivery(*Delivery) error
}
