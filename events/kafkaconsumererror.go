package events

import (
	"errors"
	"fmt"

	"github.com/IBM/sarama"
)

var ErrRegisterHanlderNotPermited = errors.New("cannot register event handler while consumer is running")

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
