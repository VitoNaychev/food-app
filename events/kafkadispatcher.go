package events

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type KafkaEventDispatcher struct {
	producer sarama.SyncProducer
}

func NewKafkaEventDispatcher(brokersAddrs []string) (*KafkaEventDispatcher, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(brokersAddrs, config)
	if err != nil {
		return nil, err
	}

	return &KafkaEventDispatcher{producer: producer}, nil
}

func (k *KafkaEventDispatcher) Close() {
	k.producer.Close()
}

func (k *KafkaEventDispatcher) Publish(topic string, event interface{}) error {
	eventJSON, _ := json.Marshal(event)

	message := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(eventJSON)}
	_, _, err := k.producer.SendMessage(message)

	return err
}
