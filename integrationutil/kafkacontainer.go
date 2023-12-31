package integrationutil

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

func SetupKafkaContainer(t testing.TB) (string, []string) {
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

	containerID := kafkaContainer.GetContainerID()
	brokersAddrs, _ := kafkaContainer.Brokers(ctx)
	return containerID, brokersAddrs
}
