package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/snisid/platform/internal/platform/events"
	"github.com/snisid/platform/internal/platform/logger"
)

type AgentTask struct {
	ID        string `json:"id"`
	AgentType string `json:"agentType"`
	Command   string `json:"command"`
}

func main() {
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	
	consumer := events.NewConsumer([]string{broker}, "swarm-ir-group", "swarm.tasks")
	
	logger.Info("Incident Responder Agent started...")

	err := consumer.Read(context.Background(), func(msg kafka.Message) error {
		var task AgentTask
		json.Unmarshal(msg.Value, &task)

		if task.AgentType == "INCIDENT_RESPONDER" {
			logger.Info("AGENT-IR: Executing response playbook: " + task.Command)
			// Trigger quarantine or lockout logic...
		}

		return nil
	})

	if err != nil {
		logger.Fatal("Agent error", err)
	}
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
