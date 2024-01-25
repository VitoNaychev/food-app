package models

import "time"

type Ticket struct {
	ID           int
	State        TicketState
	RestaurantID int `db:"restaurant_id"`
	Total        float32
	ReadyBy      time.Time `db:"ready_by"`
	// PreparingTime      time.Time
	// PickedUpTime       time.Time
	// ReadyForPickupTime time.Time
}
