package models

import "time"

type Delivery struct {
	ID                int
	PickupAddressID   int
	DeliveryAddressID int
	ReadyBy           time.Time
	Status            DeliveryStatus
}
