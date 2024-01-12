package svcevents

import (
	"time"

	"github.com/VitoNaychev/food-app/events"
)

const KITCHEN_EVENTS_TOPIC = "ktichen-events-topic"

const (
	TICKET_BEGIN_PREPARING_EVENT_ID events.EventID = iota
	TICKET_FINISH_PREPARING_EVENT_ID
	TICKET_CANCEL_EVENT_ID
)

type TicketBeginPreparingEvent struct {
	ID      int
	ReadyBy time.Time
}

type TicketFinishPreparingEvent struct {
	ID int
}

type TicketCancelEvent struct {
	ID int
}
