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

func (s *TaskService) CreateTask(userID int, content string) (models.Task, error) {
	if content == "" {
		return models.Task{}, ErrContentRequired
	}
	if len(content) > 100 {
		return models.Task{}, ErrContentTooLong
	}

	return s.store.CreateTask(userID, content)

}

func (s *TaskService) GetTask(userID, taskID int) (models.Task, error) {
	task, err := s.store.GetTask(userID, taskID)
	if err != nil {
		return models.Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (s *TaskService) ListTasks(userID, limit, offset int) ([]models.Task, error) {
	task, err := s.store.ListTasks(userID, limit, offset)
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

func (s *TaskService) UpdateTask(content string, completed bool, taskID, userID int) (models.Task, error) {
	task, err := s.store.UpdateTask(userID, taskID, content, completed)
	if err != nil {
		return models.Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (s *TaskService) DeleteTask(userID, taskID int) error {
	err := s.store.DeleteTask(userID, taskID)
	if err != nil {
		return ErrTaskNotFound
	}
	return nil

}
