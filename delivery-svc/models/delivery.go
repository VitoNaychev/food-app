package models

import "time"

type Delivery struct {
	ID                int
	CourierID         int
	PickupAddressID   int
	DeliveryAddressID int
	ReadyBy           time.Time
	State             DeliveryState
}
