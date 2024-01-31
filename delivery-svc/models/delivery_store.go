package models

type DeliveryStore interface {
	CreateDelivery(*Delivery) error
	GetDeliveryByID(int) (Delivery, error)
	UpdateDelivery(*Delivery) error
	GetActiveDeliveryByCourierID(courierID int) (Delivery, error)
}
