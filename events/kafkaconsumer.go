package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
)

type ConsumerError struct {
	Topic     string
	Partition int32
	Err       error
}

func (c *ConsumerError) Error() string {
	return fmt.Sprintf("error while consuming %s/%d: %v", c.Topic, c.Partition, c.Err)
}

func (c *ConsumerError) Unwrap() error {
	return c.Err
}

func SaramaConsumerErrorToEventsConsumerError(saramaErr sarama.ConsumerError) ConsumerError {
	err := ConsumerError{
		Topic:     saramaErr.Topic,
		Partition: saramaErr.Partition,
		Err:       saramaErr.Err,
	}

	return err
}

type BaseKafkaEventHandler struct {
	EventHandler EventHandlerFunc
}

func NewKafkaEventHandler(eventHandler EventHandlerFunc) BaseKafkaEventHandler {
	kafkeEventHandler := BaseKafkaEventHandler{}
	kafkeEventHandler.EventHandler = eventHandler

	return kafkeEventHandler
}

func (b *BaseKafkaEventHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (b *BaseKafkaEventHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claims sarama.ConsumerGroupClaim) error {
	for message := range claims.Messages() {
		var unmarshalEvent UnmarshalEvent
		json.Unmarshal(message.Value, &unmarshalEvent)

		envelope := EventEnvelope{
			EventID:     unmarshalEvent.EventID,
			AggregateID: unmarshalEvent.AggregateID,
			Timestamp:   unmarshalEvent.Timestamp,
		}

		err := b.EventHandler(envelope, unmarshalEvent.Payload)
		sess.MarkMessage(message, "")

		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BaseKafkaEventHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

type KafkaEventConsumer struct {
	ctx        context.Context
	cancel     context.CancelFunc
	group      sarama.ConsumerGroup
	handlersWg sync.WaitGroup

	ErrorsChan chan error
}

func NewKafkaEventConsumer(brokersAddrs []string, groupID string) (*KafkaEventConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup(brokersAddrs, groupID, config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	kafkaEventConsumer := KafkaEventConsumer{
		ctx:        ctx,
		cancel:     cancel,
		group:      group,
		handlersWg: sync.WaitGroup{},

		ErrorsChan: make(chan error),
	}

	return &kafkaEventConsumer, nil
}

func (k *KafkaEventConsumer) Close() {
	k.cancel()
	k.handlersWg.Wait()
	k.group.Close()
}

func (k *KafkaEventConsumer) RegisterEventHandler(topic string, eventHandler EventHandlerFunc) {
	consumerGroupHandler := NewKafkaEventHandler(eventHandler)

	k.handlersWg.Add(1)
	go k.handleEvents(topic, &consumerGroupHandler)
}

func (k *KafkaEventConsumer) handleEvents(topic string, consumerGroupHandler sarama.ConsumerGroupHandler) {
	for {
		select {
		case <-k.ctx.Done():
			k.handlersWg.Done()
			return
		case err := <-k.group.Errors():
			saramaErr := *err.(*sarama.ConsumerError)
			eventsErr := SaramaConsumerErrorToEventsConsumerError(saramaErr)

			k.ErrorsChan <- &eventsErr
		default:
			err := k.group.Consume(k.ctx, []string{topic}, consumerGroupHandler)
			if err != nil {

			}
		}
	}
}
