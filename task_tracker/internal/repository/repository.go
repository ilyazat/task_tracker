package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/ilyazat/task_tracker/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, task model.Task) error
	Close(ctx context.Context, taskID uuid.UUID) error
	Open(ctx context.Context, taskID uuid.UUID) error
}
