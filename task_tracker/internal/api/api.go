package api

import (
	"asyncArchCourse/task_tracker/internal/model"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type taskService interface {
	CreateTask(ctx context.Context, task model.Task) error
	OpenTask(ctx context.Context, taskID uuid.UUID) error
	CloseTask(ctx context.Context, taskID uuid.UUID) error
}

type TaskHandler struct {
	service taskService
}

func NewTaskHandler(service taskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (ts *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	if err := ts.service.CreateTask(ctx, task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ts *TaskHandler) OpenTask(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	taskID, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/tasks/"))
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := ts.service.OpenTask(ctx, taskID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ts *TaskHandler) CloseTask(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	taskID, err := uuid.Parse(strings.TrimPrefix(r.URL.Path, "/tasks/"))
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := ts.service.CloseTask(ctx, taskID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
