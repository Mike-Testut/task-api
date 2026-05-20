package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/mike-testut/task-api/internal/handlers"
	"github.com/mike-testut/task-api/internal/service"
	"github.com/mike-testut/task-api/internal/store"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	
)

func connectToDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn, ok := os.LookupEnv("DB_DSN")
	if !ok {
		dsn = "postgres://postgres:" + os.Getenv("POSTGRES_PW") + "@localhost:5432/postgres?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}
	log.Println("Successfully connected to database")
	return db

}
func main() {
	db := connectToDB()
	defer db.Close()


	taskStore := store.NewPostgresStore(db)
	userStore := store.NewPostgresUserStore(db)

	taskService := service.NewTaskService(taskStore)
	authService := service.NewAuthService(userStore)

	taskHandlers := handlers.NewTaskHandlers(taskService)
	authHandlers := handlers.NewAuthHandlers(authService)

	authMiddleware := handlers.NewAuthMiddleware(authService)
	mainRouter := http.NewServeMux()
	protectedRouter := http.NewServeMux()
	

	protectedRouter.HandleFunc("GET /tasks", taskHandlers.ListTasksHandler)
	protectedRouter.HandleFunc("POST /tasks", taskHandlers.CreateTaskHandler)
	protectedRouter.HandleFunc("GET /tasks/{id}", taskHandlers.GetTaskHandler)
	protectedRouter.HandleFunc("PUT /tasks/{id}", taskHandlers.UpdateTaskHandler)
	protectedRouter.HandleFunc("DELETE /tasks/{id}", taskHandlers.DeleteTaskHandler)

	mainRouter.HandleFunc("POST /v1/register", authHandlers.RegisterHandler)
	mainRouter.HandleFunc("POST /v1/login", authHandlers.LoginHandler)

	mainRouter.Handle("/v1/", http.StripPrefix("/v1", authMiddleware.Protect(protectedRouter)))
	

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	addr := ":" + port

	loggingHandler := handlers.LoggingMiddleware(mainRouter)

	server := &http.Server{
		Addr:    addr,
		Handler: loggingHandler,
	}
	

	log.Printf("Starting server on %s...", addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
