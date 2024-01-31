package stubs

import "github.com/VitoNaychev/food-app/events"

type StubEventPublisher struct {
	SpyTopic string
	SpyEvent events.InterfaceEvent
}

func (s *StubEventPublisher) Publish(topic string, event events.InterfaceEvent) error {
	s.SpyTopic = topic
	s.SpyEvent = event

	return nil
}
