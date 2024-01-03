package events_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/testutil"
)

var DummyError = errors.New("dummy error")

type SpyHandler struct {
	message string
	err     error
}

func (s *SpyHandler) EventHandler(event events.EventEnvelope, payload []byte) error {
	var dummyEvent DummyEvent
	json.Unmarshal(payload, &dummyEvent)
	s.message = dummyEvent.Message

	return s.err
}

func TestKafkaEventConsumer(t *testing.T) {
	containerID, brokersAddrs := integrationutil.SetupKafkaContainer(t)

	kafkaEventConsumer, err := events.NewKafkaEventConsumer(brokersAddrs, "test-group")
	testutil.AssertNoErr(t, err)
	t.Cleanup(kafkaEventConsumer.Close)

	spy := SpyHandler{}

	topic := "test-topic"
	kafkaEventConsumer.RegisterEventHandler(topic, spy.EventHandler)

	t.Run("receives a message", func(t *testing.T) {
		payload := DummyEvent{"Hello, World"}
		event := events.NewEvent(1, 1, payload)

		message, _ := json.Marshal(event)
		produceMessage(t, containerID, topic, string(message))

		testutil.AssertEqual(t, spy.message, payload.Message)
	})

	t.Run("writes error from event handler to ErrorsChan", func(t *testing.T) {
		spy.err = DummyError

		payload := DummyEvent{"Hello, World"}
		event := events.NewEvent(1, 1, payload)

		message, _ := json.Marshal(event)
		produceMessage(t, containerID, topic, string(message))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		t.Cleanup(cancel)

		var err error
		select {
		case err = <-kafkaEventConsumer.ErrorsChan:
			consumerError := err.(*events.ConsumerError)
			testutil.AssertEqual(t, consumerError.Err, spy.err)
		case <-ctx.Done():
			t.Fatalf("didn't receive error before timeout")
		}
	})
}
