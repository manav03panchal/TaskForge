// cmd/api/main.go
package main

import (
	"log"
	"net"

	v1 "github.com/manav03panchal/taskforge/internal/api/proto/v1"
	"github.com/manav03panchal/taskforge/internal/kafka"
	"github.com/manav03panchal/taskforge/internal/redis"
	"github.com/manav03panchal/taskforge/internal/service"
	"google.golang.org/grpc"
)

func main() {
	// Initialize Kafka producer
	kafkaProducer, err := kafka.NewProducer([]string{"localhost:9092"}, "tasks")
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Initialize Redis client
	redisClient, err := redis.NewClient("localhost:6379")
	if err != nil {
		log.Fatalf("failed to create redis client: %v", err)
	}

	// Create gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	taskService := service.NewTaskService(kafkaProducer, redisClient)
	v1.RegisterTaskServiceServer(s, taskService)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
