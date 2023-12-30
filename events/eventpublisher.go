package events

type EventPublisher interface {
	Publish(string, Event) error
}
