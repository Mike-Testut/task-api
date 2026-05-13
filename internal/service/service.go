package service

import (
	"fmt"

	"github.com/mike-testut/task-api/internal/models"
	"github.com/mike-testut/task-api/internal/store"
)

type TaskService struct {
	store *store.TaskStore
}

func NewTaskService(s *store.TaskStore) *TaskService {
	return &TaskService{store: s}
}

func (s *TaskService) CreateTask(content string) (models.Task, error) {
	if content == "" {
		return models.Task{}, fmt.Errorf("task content cannot be empty")
	}
	if len(content) > 100 {
		return models.Task{}, fmt.Errorf("task content is too long (max 100 characters)")
	}

	task := s.store.CreateTask(content)

	return task, nil
}

func (s *TaskService) GetTask(id int)(models.Task, error){
	return s.store.GetTask(id)
}

func (s *TaskService) ListTasks() []models.Task{
	return s.store.ListTasks()
}
