package events

import (
	"context"
	"errors"
	"reflect"

	"github.com/IBM/sarama"
)

type KafkaEventConsumer struct {
	group                sarama.ConsumerGroup
	groupHandler         KafkaConsumerGroupHandler
	topicRegistry        TopicRegistry
	eventHandlerRegistry map[RegistryKey]RegistryEntry

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

	kafkaEventConsumer := KafkaEventConsumer{
		group: group,
		groupHandler: KafkaConsumerGroupHandler{
			eventHandlerRegistry: map[RegistryKey]RegistryEntry{},
		},
		topicRegistry:        TopicRegistry{},
		eventHandlerRegistry: map[RegistryKey]RegistryEntry{},

		ErrorsChan: make(chan error),
	}

	return &kafkaEventConsumer, nil
}

func (k *KafkaEventConsumer) RegisterEventHandler(topic string, eventID EventID, eventHandler InterfaceEventHandler, eventType reflect.Type) error {
	k.topicRegistry[topic] = true

	entry := RegistryEntry{
		eventHandler: eventHandler,
		eventType:    eventType,
	}
	k.eventHandlerRegistry[GetRegistryKey(topic, eventID)] = entry

	return nil
}

func (k *KafkaEventConsumer) Run(ctx context.Context) {
	topics := k.topicRegistry.GetTopics()
	k.groupHandler.eventHandlerRegistry = k.eventHandlerRegistry

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-k.group.Errors():
			if err != nil {
				k.handleConsumerError(err)
			}
		default:
			err := k.group.Consume(ctx, topics, &k.groupHandler)
			if err != nil {
				// figure out what to do when group.Consume returns an error
			}
		}
	}
}

func (k *KafkaEventConsumer) Close() {
	k.group.Close()
}

func (k *KafkaEventConsumer) handleConsumerError(err error) {
	saramaErr := &sarama.ConsumerError{}
	if errors.As(err, &saramaErr) {
		eventsErr := SaramaConsumerErrorToEventsConsumerError(*saramaErr)
		k.ErrorsChan <- &eventsErr
	} else {
		k.ErrorsChan <- err
	}
}
