package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/mike-testut/task-api/internal/handlers"
	"github.com/mike-testut/task-api/internal/router"
	"github.com/mike-testut/task-api/internal/service"
	"github.com/mike-testut/task-api/internal/store"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
	authHandler := handlers.NewAuthHandlers(authService)

	
	mainRouter := http.NewServeMux()
	
	v1Router := http.NewServeMux()

	v1Router.HandleFunc("GET /tasks", taskHandlers.ListTasksHandler)
	v1Router.HandleFunc("POST /tasks", taskHandlers.CreateTaskHandler)
	v1Router.HandleFunc("GET /tasks/{id}", taskHandlers.GetTaskHandler)
	v1Router.HandleFunc("PUT /tasks/{id}", taskHandlers.UpdateTaskHandler)
	v1Router.HandleFunc("DELETE /tasks/{id}", taskHandlers.DeleteTaskHandler)

	v1Router.HandleFunc("POST /register", authHandler.RegisterHandler)
	v1Router.HandleFunc("POST /login", authHandler.LoginHandler)

	mainRouter.Handle("/v1/", http.StripPrefix("/v1", v1Router))
	

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
