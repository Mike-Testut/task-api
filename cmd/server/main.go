package main

import (
	"log"
	"net/http"

	"github.com/mike-testut/task-api/internal/handlers"
	"github.com/mike-testut/task-api/internal/router"
	"github.com/mike-testut/task-api/internal/service"
	"github.com/mike-testut/task-api/internal/store"
)

func main() {
	taskStore := store.NewTaskStore()

	taskService := service.NewTaskService(taskStore)

	taskHandlers := handlers.NewTaskHandlers(taskService)

	appRouter := router.NewRouter(taskHandlers)

	server := &http.Server{
		Addr:    ":8080",
		Handler: appRouter,
	}

	log.Println("Starting server on localhost:8080...")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
