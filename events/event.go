package events

import (
	"encoding/json"
	"time"
)

type EventID int

type EventEnvelope struct {
	EventID     EventID
	AggregateID int
	Timestamp   time.Time
}

func NewEventEnvelope(eventID EventID, aggregateID int) EventEnvelope {
	return EventEnvelope{
		EventID:     eventID,
		AggregateID: aggregateID,
		Timestamp:   time.Now().Round(0),
	}
}

type Event struct {
	EventID     EventID
	AggregateID int
	Timestamp   time.Time
	Payload     interface{}
}

type UnmarshalEvent struct {
	EventID     EventID
	AggregateID int
	Timestamp   time.Time
	Payload     json.RawMessage
}

type GenericEvent[T any] struct {
	EventID     EventID
	AggregateID int
	Timestamp   time.Time
	Payload     T
}

func NewEvent(eventID EventID, aggregateID int, payload interface{}) Event {
	return Event{
		EventID:     eventID,
		AggregateID: aggregateID,
		Timestamp:   time.Now().Round(0),
		Payload:     payload,
	}
}
