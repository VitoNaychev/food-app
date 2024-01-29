package models

import "time"

type Delivery struct {
	ID                int
	CourierID         int       `db:"corier_id"`
	PickupAddressID   int       `db:"pickup_address_id"`
	DeliveryAddressID int       `db:"delivery_address_id"`
	ReadyBy           time.Time `db:"ready_by"`
	State             DeliveryState
}
