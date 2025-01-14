package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	v1 "github.com/manav03panchal/taskforge/internal/api/proto/v1"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	// Add these configurations
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 100 * time.Millisecond

	// Add better client ID
	config.ClientID = "taskforge-producer"

	// Debug logging
	log.Printf("Connecting to Kafka brokers: %v", brokers)

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *Producer) SendTask(task *v1.Task) error {
	// Convert task to JSON
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(task.Id),
		Value: sarama.ByteEncoder(taskJSON),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Task sent to Kafka - Topic: %s, Partition: %d, Offset: %d", p.topic, partition, offset)
	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
