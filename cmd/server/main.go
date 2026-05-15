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

	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	taskService := service.NewTaskService(taskStore)

	taskHandlers := handlers.NewTaskHandlers(taskService)

	appRouter := router.NewRouter(taskHandlers)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	addr := ":" + port

	server := &http.Server{
		Addr:    addr,
		Handler: appRouter,
	}

	log.Printf("Starting server on %s...", addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
