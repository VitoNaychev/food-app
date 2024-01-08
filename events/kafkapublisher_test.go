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
	event := UnmarshalEvent[DummyEvent]([]byte(message))
	got := EventToInterfaceEvent(event)

	testutil.AssertEqual(t, got, want)
}

func UnmarshalEvent[T any](data []byte) (event events.Event[T]) {
	json.Unmarshal(data, &event)
	return
}

func EventToInterfaceEvent[T any](generic events.Event[T]) events.InterfaceEvent {
	event := events.InterfaceEvent{
		EventID:     generic.EventID,
		AggregateID: generic.AggregateID,
		Timestamp:   generic.Timestamp,
		Payload:     generic.Payload,
	}

	return event
}
