package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	v1 "github.com/manav03panchal/taskforge/internal/api/proto/v1"
	"github.com/manav03panchal/taskforge/internal/interfaces"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskService struct {
	v1.UnimplementedTaskServiceServer
	kafkaProducer interfaces.KafkaProducer
	redisClient   interfaces.RedisClient
}

func NewTaskService(kafkaProducer interfaces.KafkaProducer, redisClient interfaces.RedisClient) *TaskService {
	return &TaskService{
		kafkaProducer: kafkaProducer,
		redisClient:   redisClient,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, req *v1.CreateTaskRequest) (*v1.Task, error) {
	task := &v1.Task{
		Id:        uuid.New().String(),
		Type:      req.Type,
		Command:   req.Command,
		Status:    "PENDING",
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	log.Printf("Creating task: %v", task.Id)
	// TODO: Save task to redis !DONE
	err := s.redisClient.SaveTask(ctx, task)
	if err != nil {
		log.Fatalf("failed to save task to redis: %v", err)
		return nil, err
	}
	log.Printf("Task saved to redis: %v", task.Id)
	// TODO: Send task to Kafka !DONE
	err = s.kafkaProducer.SendTask(task)
	if err != nil {
		log.Printf("Failed to send task to Kafka: %v", err)
		return nil, err
	}
	log.Printf("Task sent to Kafka: %v", task.Id)

	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, req *v1.GetTaskRequest) (*v1.Task, error) {
	// TODO: Get from Redis
	log.Printf("Getting task: %v", req.Id)

	task, err := s.redisClient.GetTask(ctx, req.Id)
	if err != nil {
		log.Printf("Failed to get task from Redis: %v", err)
		return nil, err
	}
	return task, nil
}

func (s *TaskService) ListTasks(ctx context.Context, req *v1.ListTasksRequest) (*v1.ListTasksResponse, error) {
	// TODO: List from Redis
	return nil, nil
}
