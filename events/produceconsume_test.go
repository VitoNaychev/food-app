package events_test

import (
	"testing"

	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestProduceConsume(t *testing.T) {
	containerID, _ := integrationutil.SetupKafkaContainer(t)

	t.Run("produce and consume message from topic-A", func(t *testing.T) {
		topic := "topic-A"
		want := "Hello, World"

		produceMessage(t, containerID, topic, want)
		got := consumeMessage(t, containerID, topic)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("produce and consume message from topic-B", func(t *testing.T) {
		topic := "topic-B"
		want := "Hello, World"

		produceMessage(t, containerID, topic, want)
		got := consumeMessage(t, containerID, topic)

		testutil.AssertEqual(t, got, want)
	})
}
