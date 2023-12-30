package events

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/IBM/sarama"
)

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
		if err != nil {
			return err
		}

		sess.MarkMessage(message, "")
	}
	return nil
}

func (b *BaseKafkaEventHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

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
	go k.handleEvents(topic, consumerGroupHandler)
}

func (k *KafkaEventConsumer) handleEvents(topic string, consumerGroupHandler sarama.ConsumerGroupHandler) {
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
