package main

import (
	"asyncArchCourse/task_tracker/internal/api"
	"asyncArchCourse/task_tracker/internal/repository/simpletasks"
	"asyncArchCourse/task_tracker/internal/service/tasks"
	"log"
	"net/http"
)

func main() {

	storage := simpletasks.NewStorage()

	taskSvc := tasks.NewTaskService(storage)

	handler := api.NewTaskHandler(taskSvc)

	http.HandleFunc("/tasks", handler.CreateTask)
	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			handler.OpenTask(w, r)
		case http.MethodPatch:
			handler.CloseTask(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
