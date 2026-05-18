package service

import (
	"errors"

	"github.com/mike-testut/task-api/internal/models"
	"github.com/mike-testut/task-api/internal/store"
)

var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrContentRequired = errors.New("task content is required")
	ErrContentTooLong  = errors.New("task content is too long(max 100 characters)")
)

type TaskService struct {
	store store.Store
}

func NewTaskService(s store.Store) *TaskService {
	return &TaskService{store: s}
}

func (s *TaskService) CreateTask(content string) (models.Task, error) {
	if content == "" {
		return models.Task{}, ErrContentRequired
	}
	if len(content) > 100 {
		return models.Task{}, ErrContentTooLong
	}

	return s.store.CreateTask(content)

}

func (s *TaskService) GetTask(id int) (models.Task, error) {
	task, err := s.store.GetTask(id)
	if err != nil {
		return models.Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (s *TaskService) ListTasks(limit, offset int) ([]models.Task, error) {
	task, err := s.store.ListTasks(limit, offset)
	if err != nil {
		return []models.Task{}, ErrTaskNotFound
	}
	if limit <= 0 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return task, nil

}

func (s *TaskService) UpdateTask(id int, content string, completed bool) (models.Task, error) {
	task, err := s.store.UpdateTask(id, content, completed)
	if err != nil {
		return models.Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (s *TaskService) DeleteTask(id int) error {
	err := s.store.DeleteTask(id)
	if err != nil {
		return ErrTaskNotFound
	}
	return nil

}
