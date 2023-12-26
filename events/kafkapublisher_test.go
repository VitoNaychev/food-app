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
	event := DummyEvent{"Hello, World!"}

	err = publisher.Publish(topic, event)
	testutil.AssertNil(t, err)

	message := consumeMessage(t, containerID, topic)

	var got DummyEvent
	json.Unmarshal([]byte(message), &got)
	testutil.AssertEqual(t, got, event)
}
