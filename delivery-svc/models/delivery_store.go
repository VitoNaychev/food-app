package models

type DeliveryStore interface {
	GetDeliveryByID(int) (Delivery, error)
	UpdateDelivery(*Delivery) error
	GetActiveDeliveryByCourierID(courierID int) (Delivery, error)
}
