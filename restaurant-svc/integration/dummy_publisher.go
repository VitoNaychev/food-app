package integration

import "github.com/VitoNaychev/food-app/events"

type DummyPublisher struct {
}

func (s *DummyPublisher) Publish(topic string, event events.InterfaceEvent) error {
	return nil
}
