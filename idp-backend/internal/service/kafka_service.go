package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const DeployRequestsTopic = "deploy-requests"

type KafkaService struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaService(brokers string) *KafkaService {
	// Initialize Producer (Writer)
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    DeployRequestsTopic,
		Balancer: &kafka.LeastBytes{},
	}

	// Initialize Consumer (Reader)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokers},
		GroupID:  "idp-orchestrator",
		Topic:    DeployRequestsTopic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	log.Println("Kafka Service initialized with segmentio/kafka-go")
	return &KafkaService{
		writer: w,
		reader: r,
	}
}

func (s *KafkaService) ProduceDeployRequest(projectID uuid.UUID) error {
	msg := map[string]string{"project_id": projectID.String()}
	value, _ := json.Marshal(msg)

	err := s.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(projectID.String()),
			Value: value,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %v", err)
	}

	return nil
}

func (s *KafkaService) StartConsumer(handler func(projectID uuid.UUID) error) {
	go func() {
		for {
			m, err := s.reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Kafka reader error: %v", err)
				time.Sleep(time.Second) // Prevent tight loop on error
				continue
			}

			var data map[string]string
			if err := json.Unmarshal(m.Value, &data); err != nil {
				log.Printf("Failed to unmarshal Kafka message: %v", err)
				continue
			}

			projectID, err := uuid.Parse(data["project_id"])
			if err != nil {
				log.Printf("[KAFKA] Invalid project ID in message: %v", err)
				continue
			}

			log.Printf("[KAFKA] RECEIVED: deploy request for project %s", projectID)
			
			// Process message
			start := time.Now()
			if err := handler(projectID); err != nil {
				log.Printf("[KAFKA] ERROR: Failed to process deployment for %s: %v", projectID, err)
			} else {
				log.Printf("[KAFKA] SUCCESS: Deployment for %s completed in %v", projectID, time.Since(start))
			}
		}
	}()
}

func (s *KafkaService) Close() {
	if err := s.writer.Close(); err != nil {
		log.Printf("Error closing Kafka writer: %v", err)
	}
	if err := s.reader.Close(); err != nil {
		log.Printf("Error closing Kafka reader: %v", err)
	}
}
