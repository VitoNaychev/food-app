package events

import (
	"context"
	"log"
)

func LogEventConsumerErrors(ctx context.Context, eventConsumer *KafkaEventConsumer) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-eventConsumer.ErrorsChan:
			log.Println("Kafka Event Consumer: ", err)
		}
	}
}
