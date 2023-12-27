package events_test

import (
	"encoding/json"
	"testing"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestKafkaPublisher(t *testing.T) {
	containerID, brokersAddrs := SetupKafkaContainer(t)

	publisher, err := events.NewKafkaEventPublisher(brokersAddrs)
	testutil.AssertNil(t, err)
	defer publisher.Close()

	topic := "test-topic"
	payload := DummyEvent{"Hello, World!"}

	event := events.NewEvent(1, 1, payload)

	err = publisher.Publish(topic, event)
	testutil.AssertNil(t, err)

	message := consumeMessage(t, containerID, topic)

	var got events.GenericEvent
	json.Unmarshal([]byte(message), &got)

	var gotPayload DummyEvent
	json.Unmarshal(got.Payload, &gotPayload)

	testutil.AssertEqual(t, gotPayload, payload)
}
