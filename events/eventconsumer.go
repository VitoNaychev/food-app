package events

type EventHandlerFunc func(interface{}) error

type EventConsumer interface {
	RegisterEventHandler(string, EventHandlerFunc)
}
