package router

import (
	"net/http"

	"github.com/mike-testut/task-api/internal/handlers"
)

func NewRouter(th *handlers.TaskHandlers) http.Handler{
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", th.ListTasksHandler)
	mux.HandleFunc("POST /tasks", th.CreateTaskHandler)
	mux.HandleFunc("GET /tasks/{id}", th.GetTaskHandler)
	mux.HandleFunc("PUT /tasks/{id}", th.UpdateTaskHandler)
	mux.HandleFunc("DELETE /tasks/{id}", th.DeleteTaskHandler)
	

	var wrappedMux http.Handler = mux
	wrappedMux = handlers.LoggingMiddleware(wrappedMux)

	return wrappedMux
}
