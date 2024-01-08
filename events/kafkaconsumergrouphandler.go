package events

import (
	"encoding/json"
	"reflect"

	"github.com/IBM/sarama"
)

type KafkaConsumerGroupHandler struct {
	eventHandlerRegistry map[RegistryKey]RegistryEntry
}

func (b *KafkaConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (b *KafkaConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claims sarama.ConsumerGroupClaim) error {
	for message := range claims.Messages() {
		var rawPayloadEvent RawPayloadEvent
		json.Unmarshal(message.Value, &rawPayloadEvent)

		if registryEntry, ok := b.eventHandlerRegistry[GetRegistryKey(claims.Topic(), rawPayloadEvent.EventID)]; ok {
			payloadPtr := reflect.New(registryEntry.eventType).Interface()
			json.Unmarshal(rawPayloadEvent.Payload, payloadPtr)

			payload := reflect.ValueOf(payloadPtr).Elem().Interface()
			event := InterfaceEvent{
				EventID:     rawPayloadEvent.EventID,
				AggregateID: rawPayloadEvent.AggregateID,
				Timestamp:   rawPayloadEvent.Timestamp,
				Payload:     payload,
			}

			err := registryEntry.eventHandler(event)
			if err != nil {
				return err
			}
		}
		sess.MarkMessage(message, "")
	}
	return nil
}

func (b *KafkaConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
