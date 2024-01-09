package models

import "time"

type Ticket struct {
	ID                 int
	State              TicketState
	RestaurantID       int
	ReadyBy            time.Time
	AcceptTime         time.Time
	PreparingTime      time.Time
	PickedUpTime       time.Time
	ReadyForPickupTime time.Time
}
