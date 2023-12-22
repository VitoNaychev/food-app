package events_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/IBM/sarama"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

type DummyEvent struct {
	Message string
}

func SetupKafkaContainer(t testing.TB) []string {
	ctx := context.Background()

	kafkaContainer, err := kafka.RunContainer(ctx,
		kafka.WithClusterID("test-cluster"),
		testcontainers.WithImage("confluentinc/confluent-local:7.5.0"),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := kafkaContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	})

	brokersAddrs, _ := kafkaContainer.Brokers(ctx)
	return brokersAddrs
}

func TestKafkaDispatcher(t *testing.T) {
	brokersAddrs := SetupKafkaContainer(t)

	dispatcher, err := events.NewKafkaEventDispatcher(brokersAddrs)
	testutil.AssertNil(t, err)
	defer dispatcher.Close()

	topic := "test-topic"
	event := DummyEvent{"Hello, World!"}

	err = dispatcher.Publish(topic, event)
	testutil.AssertNil(t, err)

	consumer, err := sarama.NewConsumer(brokersAddrs, sarama.NewConfig())
	if err != nil {
		t.Fatal("NewConsumer err: ", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		t.Fatal("ConsumePartition err: ", err)
	}
	defer partitionConsumer.Close()

	message := <-partitionConsumer.Messages()

	got := DummyEvent{}
	json.Unmarshal(message.Value, &got)

	testutil.AssertEqual(t, got, event)
}
