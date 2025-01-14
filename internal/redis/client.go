// internal/redis/client.go
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	v1 "github.com/manav03panchal/taskforge/internal/api/proto/v1"
)

type Client struct {
	client *redis.Client
}

func NewClient(addr string) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) SaveTask(ctx context.Context, task *v1.Task) error {
	key := fmt.Sprintf("task:%s", task.Id)

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	err = c.client.Set(ctx, key, taskJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to save to redis: %v", err)
	}

	log.Printf("Task saved to Redis with key: %s", key)
	return nil
}

func (c *Client) GetTask(ctx context.Context, id string) (*v1.Task, error) {
	key := fmt.Sprintf("task:%s", id)

	// Get from Redis
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get from redis: %v", err)
	}

	// Parse JSON
	var task v1.Task
	err = json.Unmarshal([]byte(val), &task)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %v", err)
	}

	log.Printf("Task retrieved from Redis: %s", id)
	return &task, nil
}
