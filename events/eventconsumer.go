package events

type EventHandlerFunc func(envelope EventEnvelope, payload []byte) error

type EventConsumer interface {
	RegisterEventHandler(string, EventHandlerFunc)
}
