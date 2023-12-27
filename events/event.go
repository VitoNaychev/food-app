package events

import (
	"encoding/json"
	"time"
)

type EventEnvelope struct {
	EventID     int
	AggregateID int
	Timestamp   time.Time
}

type Event struct {
	EventID     int
	AggregateID int
	Timestamp   time.Time
	Payload     interface{}
}

type UnmarshalEvent struct {
	EventID     int
	AggregateID int
	Timestamp   time.Time
	Payload     json.RawMessage
}

type GenericEvent[T any] struct {
	EventID     int
	AggregateID int
	Timestamp   time.Time
	Payload     T
}

func NewEvent(eventID, aggregateID int, payload interface{}) Event {
	return Event{
		EventID:     eventID,
		AggregateID: aggregateID,
		Timestamp:   time.Now(),
		Payload:     payload,
	}
}
