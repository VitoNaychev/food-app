package stubs

import "github.com/VitoNaychev/food-app/events"

type StubEventPublisher struct {
	Topic string
	Event events.InterfaceEvent
}

func (s *StubEventPublisher) Publish(topic string, event events.InterfaceEvent) error {
	s.Topic = topic
	s.Event = event

	return nil
}
