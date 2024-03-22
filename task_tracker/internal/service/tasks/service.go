package tasks

import (
	"context"
	"github.com/google/uuid"
	"github.com/ilyazat/task_tracker/internal/model"
)

type Task struct {
	Description string
	Status      string
	Assignee    string
}

type storage interface {
	Create(ctx context.Context, task model.Task) error
	Close(ctx context.Context, taskID uuid.UUID) error
	Open(ctx context.Context, taskID uuid.UUID) error
}

type TaskService struct {
	storage storage
}

func NewTaskService(db storage) *TaskService {
	return &TaskService{
		storage: db,
	}
}

func (ts *TaskService) CreateTask(ctx context.Context, task model.Task) error {
	if err := ts.storage.Create(ctx, task); err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) OpenTask(ctx context.Context, taskID uuid.UUID) error {
	if err := ts.storage.Open(ctx, taskID); err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) CloseTask(ctx context.Context, taskID uuid.UUID) error {
	if err := ts.storage.Close(ctx, taskID); err != nil {
		return err
	}

	return nil
}
