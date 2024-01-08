package events_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
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

func (s *SpyHandler) EventHandler(event events.Event[DummyEvent]) error {
	s.message = event.Payload.Message

	return s.err
}

func TestKafkaEventConsumer(t *testing.T) {
	containerID, brokersAddrs := integrationutil.SetupKafkaContainer(t)

	kafkaEventConsumer, err := events.NewKafkaEventConsumer(brokersAddrs, "test-group")
	testutil.AssertNoErr(t, err)
	t.Cleanup(kafkaEventConsumer.Close)

	t.Run("registers an event handler and receives a message", func(t *testing.T) {
		spy := SpyHandler{}

		topic := "test-topic"
		err := kafkaEventConsumer.RegisterEventHandler(topic, 1, events.EventHandlerWrapper(spy.EventHandler), reflect.TypeOf(DummyEvent{}))
		testutil.AssertNoErr(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		go kafkaEventConsumer.Run(ctx)
		t.Cleanup(cancel)

		payload := DummyEvent{"Hello, World"}
		event := events.NewEvent(1, 1, payload)

		message, _ := json.Marshal(event)
		produceMessage(t, containerID, topic, string(message))

		testutil.AssertEqual(t, spy.message, payload.Message)
	})

	t.Run("registers two event handlers on different topics and receives messages", func(t *testing.T) {
		topicA := "topic-A"
		handlerA := SpyHandler{}
		kafkaEventConsumer.RegisterEventHandler(topicA, 1, events.EventHandlerWrapper(handlerA.EventHandler), reflect.TypeOf(DummyEvent{}))

		topicB := "topic-B"
		handlerB := SpyHandler{}
		kafkaEventConsumer.RegisterEventHandler(topicB, 1, events.EventHandlerWrapper(handlerB.EventHandler), reflect.TypeOf(DummyEvent{}))

		ctx, cancel := context.WithCancel(context.Background())
		go kafkaEventConsumer.Run(ctx)
		t.Cleanup(cancel)

		payloadA := DummyEvent{"Message A"}
		messageA := NewMarshaledEvent(1, 1, payloadA)
		produceMessage(t, containerID, topicA, string(messageA))

		payloadB := DummyEvent{"Message B"}
		messageB := NewMarshaledEvent(1, 1, payloadB)
		produceMessage(t, containerID, topicB, string(messageB))

		testutil.AssertEqual(t, handlerA.message, payloadA.Message)
		testutil.AssertEqual(t, handlerB.message, payloadB.Message)
	})

	t.Run("registers two event handlers on different eventIDs and receives messages", func(t *testing.T) {
		commonTopic := "common-topic"

		eventIDA := events.EventID(1)
		handlerA := SpyHandler{}
		kafkaEventConsumer.RegisterEventHandler(commonTopic, eventIDA, events.EventHandlerWrapper(handlerA.EventHandler), reflect.TypeOf(DummyEvent{}))

		eventIDB := events.EventID(2)
		handlerB := SpyHandler{}
		kafkaEventConsumer.RegisterEventHandler(commonTopic, eventIDB, events.EventHandlerWrapper(handlerB.EventHandler), reflect.TypeOf(DummyEvent{}))

		ctx, cancel := context.WithCancel(context.Background())
		go kafkaEventConsumer.Run(ctx)
		t.Cleanup(cancel)

		payloadA := DummyEvent{"Message A"}
		messageA := NewMarshaledEvent(eventIDA, 1, payloadA)
		produceMessage(t, containerID, commonTopic, string(messageA))

		payloadB := DummyEvent{"Message B"}
		messageB := NewMarshaledEvent(eventIDB, 1, payloadB)
		produceMessage(t, containerID, commonTopic, string(messageB))

		testutil.AssertEqual(t, handlerA.message, payloadA.Message)
		testutil.AssertEqual(t, handlerB.message, payloadB.Message)
	})

	t.Run("writes error from event handler to ErrorsChan", func(t *testing.T) {
		errTopic := "topic-err"
		errHandler := SpyHandler{err: DummyError}
		kafkaEventConsumer.RegisterEventHandler(errTopic, 1, events.EventHandlerWrapper(errHandler.EventHandler), reflect.TypeOf(DummyEvent{}))

		kafkaCtx, kafkaCancel := context.WithCancel(context.Background())
		go kafkaEventConsumer.Run(kafkaCtx)
		t.Cleanup(kafkaCancel)

		payload := DummyEvent{"Goodbye, World"}
		message := NewMarshaledEvent(1, 1, payload)
		produceMessage(t, containerID, errTopic, string(message))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		t.Cleanup(cancel)

		var err error
		select {
		case err = <-kafkaEventConsumer.ErrorsChan:
			consumerError := &events.ConsumerError{}
			if errors.As(err, &consumerError) {
				testutil.AssertEqual(t, consumerError.Err, errHandler.err)
			} else {
				t.Errorf("expected ConsumerError, got %v", reflect.TypeOf(err))
			}
		case <-ctx.Done():
			t.Fatalf("didn't receive error before timeout")
		}
	})
}

func NewMarshaledEvent(eventID events.EventID, aggregateID int, payload interface{}) []byte {
	event := events.NewEvent(eventID, aggregateID, payload)
	message, _ := json.Marshal(event)

	return message
}
