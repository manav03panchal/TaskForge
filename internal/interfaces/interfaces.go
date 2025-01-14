// internal/interfaces/interfaces.go
package interfaces

import (
	"context"

	v1 "github.com/manav03panchal/taskforge/internal/api/proto/v1"
)

type KafkaProducer interface {
	SendTask(task *v1.Task) error
	Close() error
}

type RedisClient interface {
	SaveTask(ctx context.Context, task *v1.Task) error
	GetTask(ctx context.Context, id string) (*v1.Task, error)
}
