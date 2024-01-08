package events

type EventPublisher interface {
	Publish(string, InterfaceEvent) error
}
