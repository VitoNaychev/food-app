package events_test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

const DUMMY_EVENT_ID = 10

type DummyEvent struct {
	Message string
}

func consumeMessage(t *testing.T, containerID, topic string) string {
	cmd := exec.Command("docker", "exec", containerID, "/bin/kafka-console-consumer", "--bootstrap-server", "localhost:9092", "--topic", topic, "--from-beginning", "--timeout-ms", "1000")

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Error consuming message: %v - %v", err, string(output))
	}

	return strings.TrimSpace(string(output))
}

func produceMessage(t *testing.T, containerID, topic, message string) {
	message = strings.ReplaceAll(message, `"`, `\"`)

	dockercmd := fmt.Sprintf("echo %v > input.txt", message)
	cmd := exec.Command("docker", "exec", containerID, "sh", "-c", dockercmd)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Error producing message: %v", err)
	}

	dockercmd = fmt.Sprintf("/bin/kafka-console-producer --bootstrap-server localhost:9092 --topic %v < input.txt", topic)
	cmd = exec.Command("docker", "exec", containerID, "sh", "-c", dockercmd)
	cmd.Stdin = strings.NewReader(message)

	err = cmd.Run()
	if err != nil {
		t.Fatalf("Error producing message: %v", err)
	}
}
