package integration

import "github.com/VitoNaychev/food-app/events"

type DummyPublisher struct {
}

func (s *DummyPublisher) Publish(topic string, event events.Event) error {
	return nil
}
