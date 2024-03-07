package simpletasks

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ilyazat/task_tracker/internal/model"
	"sync"
)

const InitialStorageVolume = 1000

const (
	StatusClosed = "closed"
	StatusOpen   = "open"
)

var (
	ErrNotUniqueID = errors.New("id is not unique")
	NotFound       = "id %s could not be found"
)

type Storage struct {
	mu sync.Mutex
	m  map[uuid.UUID]model.Task
}

func NewStorage() *Storage {
	m := make(map[uuid.UUID]model.Task, InitialStorageVolume)

	return &Storage{
		m: m,
	}
}

func (s *Storage) Create(_ context.Context, task model.Task) error {
	if _, ok := s.m[task.ID]; ok {
		return ErrNotUniqueID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[task.ID] = task
	return nil
}

func (s *Storage) Close(_ context.Context, taskID uuid.UUID) error {
	task, ok := s.m[taskID]
	if !ok {
		return fmt.Errorf(NotFound, taskID)
	}
	task.Status = StatusClosed

	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[taskID] = task

	return nil
}

func (s *Storage) Open(_ context.Context, taskID uuid.UUID) error {
	task, ok := s.m[taskID]
	if !ok {
		return fmt.Errorf(NotFound, taskID)
	}
	task.Status = StatusOpen

	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[taskID] = task

	return nil
}
