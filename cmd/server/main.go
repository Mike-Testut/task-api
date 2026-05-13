package main

import (
	"log"
	"net/http"
	"os"

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

	port, ok := os.LookupEnv("PORT")
	if !ok{
		port = "8080"
	}
	addr := ":" + port

	server := &http.Server{
		Addr:    addr,
		Handler: appRouter,
	}

	log.Printf("Starting server on %s...",addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
