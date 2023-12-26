package events

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type KafkaEventPublisher struct {
	producer sarama.SyncProducer
}

func NewKafkaEventPublisher(brokersAddrs []string) (*KafkaEventPublisher, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(brokersAddrs, config)
	if err != nil {
		return nil, err
	}

	return &KafkaEventPublisher{producer: producer}, nil
}

func (k *KafkaEventPublisher) Close() {
	k.producer.Close()
}

func (k *KafkaEventPublisher) Publish(topic string, event interface{}) error {
	eventJSON, _ := json.Marshal(event)

	message := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(eventJSON)}
	_, _, err := k.producer.SendMessage(message)

	return err
}
