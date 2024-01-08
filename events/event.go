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

type Event[T any] struct {
	EventID     EventID
	AggregateID int
	Timestamp   time.Time
	Payload     T
}

func NewTypedEvent[T any](eventID EventID, aggregateID int, payload T) Event[T] {
	return Event[T]{
		EventID:     eventID,
		AggregateID: aggregateID,
		Timestamp:   time.Now().Round(0),
		Payload:     payload,
	}
}

type RawPayloadEvent struct {
	EventID     EventID
	AggregateID int
	Timestamp   time.Time
	Payload     json.RawMessage
}

type InterfaceEvent struct {
	EventID     EventID
	AggregateID int
	Timestamp   time.Time
	Payload     interface{}
}

func NewEvent(eventID EventID, aggregateID int, payload any) InterfaceEvent {
	return InterfaceEvent{
		EventID:     eventID,
		AggregateID: aggregateID,
		Timestamp:   time.Now().Round(0),
		Payload:     payload,
	}
}
