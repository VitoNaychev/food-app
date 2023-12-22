package events

type EventDispatcher interface {
	Publish(string, interface{}) error
}
