package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
)

func TypedHandlerToEventHandlerFunc[T any](typedHandler func(T) error) EventHandlerFunc {
	return EventHandlerFunc(func(i interface{}) error {
		if typedParam, ok := i.(T); ok {
			return typedHandler(typedParam)
		} else {
			return fmt.Errorf("trying to call event handler with incorrect type")
		}
	})
}

type BaseKafkaEventHandler[T any] struct {
	EventHandler EventHandlerFunc
}

func NewKafkaEventHandler[T any](eventHandler func(T) error) BaseKafkaEventHandler[T] {
	kafkeEventHandler := BaseKafkaEventHandler[T]{}
	kafkeEventHandler.EventHandler = TypedHandlerToEventHandlerFunc(eventHandler)

	return kafkeEventHandler
}

func (b *BaseKafkaEventHandler[T]) Setup(sarama.ConsumerGroupSession) error { return nil }

func (b *BaseKafkaEventHandler[T]) ConsumeClaim(sess sarama.ConsumerGroupSession, claims sarama.ConsumerGroupClaim) error {
	for message := range claims.Messages() {
		var event T
		json.Unmarshal(message.Value, &event)

		err := b.EventHandler(event)
		if err != nil {
			return err
		}

		sess.MarkMessage(message, "")
	}
	return nil
}

func (b *BaseKafkaEventHandler[T]) Cleanup(sarama.ConsumerGroupSession) error { return nil }

type KafkaEventConsumer struct {
	ctx        context.Context
	cancel     context.CancelFunc
	group      sarama.ConsumerGroup
	handlersWg sync.WaitGroup
}

func NewKafkaEventConsumer(brokersAddrs []string, groupID string) (*KafkaEventConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(brokersAddrs, groupID, config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	kafkaEventConsumer := KafkaEventConsumer{
		ctx:    ctx,
		cancel: cancel,
		group:  group,
	}

	return &kafkaEventConsumer, nil
}

func (k *KafkaEventConsumer) Close() {
	k.cancel()
	k.handlersWg.Wait()
	k.group.Close()
}

func (k *KafkaEventConsumer) RegisterEventHandler(topic string, consumerGroupHandler sarama.ConsumerGroupHandler) {
	k.handlersWg.Add(1)
	go k.HandleEvents(topic, consumerGroupHandler)
}

func (k *KafkaEventConsumer) HandleEvents(topic string, consumerGroupHandler sarama.ConsumerGroupHandler) {
	for {
		select {
		case <-k.ctx.Done():
			k.handlersWg.Done()
			return
		default:
			err := k.group.Consume(k.ctx, []string{topic}, consumerGroupHandler)
			if err != nil {
				k.cancel()
				break
			}
		}
	}
}
