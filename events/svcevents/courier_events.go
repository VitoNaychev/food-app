package svcevents

import "github.com/VitoNaychev/food-app/events"

const COURIER_EVENTS_TOPIC = "courier-events-topic"

const (
	COURIER_CREATED_EVENT_ID events.EventID = iota
	COURIER_DELETED_EVENT_ID
)

type CourierCreatedEvent struct {
	ID   int
	Name string
}

type CourierDeletedEvent struct {
	ID int
}
