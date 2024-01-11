package events

import "reflect"

type EventHandlerFunc func(envelope EventEnvelope, payload []byte) error

type EventConsumer interface {
	RegisterEventHandler(string, EventID, InterfaceEventHandler, reflect.Type) error
}
