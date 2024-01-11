package models

type Ticket struct {
	ID           int
	State        TicketState
	RestaurantID int `db:"restaurant_id"`
	Total        float32
	// ReadyBy            time.Time
	// PreparingTime      time.Time
	// PickedUpTime       time.Time
	// ReadyForPickupTime time.Time
}
