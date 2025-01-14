// cmd/client/main.go
package main

import (
	"context"
	"log"
	"time"

	pb "github.com/manav03panchal/taskforge/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTaskServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a test task
	task, err := c.CreateTask(ctx, &pb.CreateTaskRequest{
		Type:    "test",
		Command: "echo 'hello world'",
	})
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	log.Printf("Created Task: %v", task)

	// Try to get the task back
	getTask, err := c.GetTask(ctx, &pb.GetTaskRequest{
		Id: task.Id,
	})
	if err != nil {
		log.Fatalf("could not get task: %v", err)
	}
	log.Printf("Retrieved Task: %v", getTask)
}
