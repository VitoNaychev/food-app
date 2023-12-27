package events

import (
	"encoding/json"
	"time"
)

type Event struct {
	EventID     int
	AggregateID int
	Timestamp   time.Time
	Payload     interface{}
}

type EventEnvelope struct {
	EventID     int
	AggregateID int
	Timestamp   time.Time
}

type GenericEvent struct {
	EventID     int
	AggregateID int
	Timestamp   time.Time
	Payload     json.RawMessage
}

func NewEvent(eventID, aggregateID int, payload interface{}) Event {
	return Event{
		EventID:     eventID,
		AggregateID: aggregateID,
		Timestamp:   time.Now(),
		Payload:     payload,
	}
}
