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

func (s *PostgresStore) CreateTask(content string) (models.Task, error) {
	var task models.Task

	query := `INSERT INTO tasks (content) VALUES ($1) RETURNING *`

	err := s.db.QueryRowx(query, content).StructScan(&task)
	if err != nil {
		return models.Task{}, fmt.Errorf("could not create task: %v", err)
	}

	return task, nil
}

func (s *PostgresStore) GetTask(id int) (models.Task, error) {
	var task models.Task
	query := `SELECT * FROM tasks WHERE id = $1`

	err := s.db.Get(&task, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, fmt.Errorf("task with id %d not found", id)
		}
		return models.Task{}, fmt.Errorf("could not get task: %v", err)
	}
	return task, nil
}

func (s *PostgresStore) ListTasks(limit, offset int) ([]models.Task, error) {
	var tasks []models.Task
	query := `SELECT * FROM tasks ORDER BY id ASC LIMIT $1 OFFSET $2`

	err := s.db.Select(&tasks, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not list tasks: %v", err)
	}
	return tasks, nil
}

func (s *PostgresStore) UpdateTask(id int, content string, completed bool) (models.Task, error) {
	var task models.Task

	query := `UPDATE tasks SET content = $1, completed = $2 WHERE id = $3 RETURNING *`

	err := s.db.QueryRowx(query, content, completed, id).StructScan(&task)
	if err != nil {
		return models.Task{}, fmt.Errorf("could not update task: ")
	}
	return task, nil
}

func (s *PostgresStore) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not delete task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}
	return nil
}
