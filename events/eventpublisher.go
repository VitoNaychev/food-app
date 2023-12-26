package events

type EventPublisher interface {
	Publish(string, interface{}) error
}
