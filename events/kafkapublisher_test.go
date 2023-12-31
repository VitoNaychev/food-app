package events_test

import (
	"encoding/json"
	"testing"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestKafkaPublisher(t *testing.T) {
	containerID, brokersAddrs := integrationutil.SetupKafkaContainer(t)

	publisher, err := events.NewKafkaEventPublisher(brokersAddrs)
	testutil.AssertNil(t, err)
	defer publisher.Close()

	aggID := 1
	payload := DummyEvent{"Hello, World!"}
	want := events.NewEvent(DUMMY_EVENT_ID, aggID, payload)

	topic := "test-topic"
	err = publisher.Publish(topic, want)
	testutil.AssertNil(t, err)

	message := consumeMessage(t, containerID, topic)
	genericEvent := UnmarshalGenericEvent[DummyEvent]([]byte(message))
	got := GenericEventToEvent(genericEvent)

	testutil.AssertEqual(t, got, want)
}

func UnmarshalGenericEvent[T any](data []byte) (event events.GenericEvent[T]) {
	json.Unmarshal(data, &event)
	return
}

func GenericEventToEvent[T any](generic events.GenericEvent[T]) events.Event {
	event := events.Event{
		EventID:     generic.EventID,
		AggregateID: generic.AggregateID,
		Timestamp:   generic.Timestamp,
		Payload:     generic.Payload,
	}

	return event
}
