package service

import (
	"fmt"

	"github.com/mike-testut/task-api/internal/models"
	"github.com/mike-testut/task-api/internal/store"
)

type TaskService struct {
	store store.Store
}

func NewTaskService(s store.Store) *TaskService {
	return &TaskService{store: s}
}

func (s *TaskService) CreateTask(content string) (models.Task, error) {
	if content == "" {
		return models.Task{}, fmt.Errorf("task content cannot be empty")
	}
	if len(content) > 100 {
		return models.Task{}, fmt.Errorf("task content is too long (max 100 characters)")
	}

	return s.store.CreateTask(content)

}

func (s *TaskService) GetTask(id int) (models.Task, error) {
	return s.store.GetTask(id)
}

func (s *TaskService) ListTasks() ([]models.Task, error) {
	return s.store.ListTasks()

}

func (s *TaskService) UpdateTask(id int, content string, completed bool) (models.Task, error) {
	return s.store.UpdateTask(id, content, completed)
}

func (s *TaskService) DeleteTask(id int) error {
	return s.store.DeleteTask(id)
}
