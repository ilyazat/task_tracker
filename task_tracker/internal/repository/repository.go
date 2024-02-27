package repository

import (
	"asyncArchCourse/task_tracker/internal/model"
	"context"
	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(ctx context.Context, task model.Task) error
	Close(ctx context.Context, taskID uuid.UUID) error
	Open(ctx context.Context, taskID uuid.UUID) error
}
