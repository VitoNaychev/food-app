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

func (d *DummyHandlerSpy) DummyHandler(event events.EventEnvelope, payload []byte) error {
	var dummyEvent DummyEvent
	json.Unmarshal(payload, &dummyEvent)
	d.message = dummyEvent.Message

	return nil
}

func TestKafkaEventConsumer(t *testing.T) {
	containerID, brokersAddrs := SetupKafkaContainer(t)

	kafkaEventConsumer, err := events.NewKafkaEventConsumer(brokersAddrs, "test-group")
	testutil.AssertNil(t, err)

	spy := DummyHandlerSpy{}
	kafkaEventHandler := events.NewKafkaEventHandler(spy.DummyHandler)

	topic := "test-topic"
	kafkaEventConsumer.RegisterEventHandler(topic, &kafkaEventHandler)

	payload := DummyEvent{"Hello, World"}
	event := events.NewEvent(1, 1, payload)

	message, _ := json.Marshal(event)
	produceMessage(t, containerID, topic, string(message))

	testutil.AssertEqual(t, spy.message, payload.Message)

	kafkaEventConsumer.Close()
}
