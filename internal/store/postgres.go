package store

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mike-testut/task-api/internal/models"
)

type PostgresStore struct {
	db *sqlx.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		db: sqlx.NewDb(db, "pgx"),
	}
}

var _ Store = (*PostgresStore)(nil)

func (s *PostgresStore) CreateTask(userID int, content string) (models.Task, error) {
	var task models.Task

	query := `INSERT INTO tasks (content, user_id) VALUES ($1,$2) RETURNING *`

	err := s.db.QueryRowx(query, content, userID).StructScan(&task)
	if err != nil {
		return models.Task{}, fmt.Errorf("could not create task: %v", err)
	}

	return task, nil
}

func (s *PostgresStore) GetTask(userID, taskID int) (models.Task, error) {
	var task models.Task
	query := `SELECT * FROM tasks WHERE id = $2 AND user_id = $1`

	err := s.db.Get(&task, query, userID, taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, fmt.Errorf("task with id %d not found", taskID)
		}
		return models.Task{}, fmt.Errorf("could not get task: %v", err)
	}
	return task, nil
}

func (s *PostgresStore) ListTasks(userID, limit, offset int) ([]models.Task, error) {
	var tasks []models.Task
	query := `SELECT * FROM tasks WHERE user_id = $1 ORDER BY id ASC LIMIT $2 OFFSET $3`

	err := s.db.Select(&tasks, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not list tasks: %v", err)
	}
	return tasks, nil
}

func (s *PostgresStore) UpdateTask(userID, taskID int, content string, completed bool) (models.Task, error) {
	var task models.Task

	query := `UPDATE tasks SET content = $3, completed = $4 WHERE id = $2 AND user_id = $1 RETURNING *`

	err := s.db.QueryRowx(query, content, completed, taskID, userID).StructScan(&task)
	if err != nil {
		return models.Task{}, fmt.Errorf("could not update task: ")
	}
	return task, nil
}

func (s *PostgresStore) DeleteTask(userID, taskID int) error {
	query := `DELETE FROM tasks WHERE id = $2 AND user_id = $1`

	result, err := s.db.Exec(query, userID, taskID)
	if err != nil {
		return fmt.Errorf("could not delete task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", taskID)
	}
	return nil
}
