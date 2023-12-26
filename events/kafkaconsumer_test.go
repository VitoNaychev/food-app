package events_test

import (
	"encoding/json"
	"testing"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/testutil"
)

type DummyHandlerSpy struct {
	message string
}

func (d *DummyHandlerSpy) DummyHandler(event DummyEvent) error {
	d.message = event.Message

	return nil
}

func TestKafkaEventConsumer(t *testing.T) {
	containerID, brokersAddrs := SetupKafkaContainer(t)

	kafkaEventConsumer, err := events.NewKafkaEventConsumer(brokersAddrs, "test-group")
	testutil.AssertNil(t, err)

	spy := DummyHandlerSpy{}
	kafkaEventHandler := events.NewKafkaEventHandler[DummyEvent](spy.DummyHandler)

	topic := "test-topic"
	kafkaEventConsumer.RegisterEventHandler(topic, &kafkaEventHandler)

	event := DummyEvent{"Hello, World"}

	message, _ := json.Marshal(event)
	produceMessage(t, containerID, topic, string(message))

	testutil.AssertEqual(t, spy.message, event.Message)

	kafkaEventConsumer.Close()
}
