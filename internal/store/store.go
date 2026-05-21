package store

import (
	"fmt"
	"sync"

	"github.com/mike-testut/task-api/internal/models"
)

type Store interface {
	CreateTask(userID int, ontent string) (models.Task, error)
	GetTask(userID int, id int) (models.Task, error)
	ListTasks(userID int, limit, offset int) ([]models.Task, error)
	UpdateTask(userID int, id int, content string, completed bool) (models.Task, error)
	DeleteTask(userID int, id int) error
}

type TaskStore struct {
	mu     sync.Mutex
	tasks  map[int]models.Task
	nextID int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *TaskStore) CreateTask(content string) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := models.Task{
		ID:        s.nextID,
		Content:   content,
		Completed: false,
	}

	s.tasks[task.ID] = task
	s.nextID++

	return task, nil
}

func (s *TaskStore) GetTask(id int) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, fmt.Errorf("task with id %d not found", id)
	}

	return task, nil
}

func (s *TaskStore) ListTasks(limit, offset int) ([]models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	allTasks := make([]models.Task, 0, len(s.tasks))

	for _, task := range s.tasks {
		allTasks = append(allTasks, task)
	}
	return allTasks, nil
}

func (s *TaskStore) UpdateTask(id int, content string, completed bool) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, fmt.Errorf("task with id %d not found", id)
	}
	task.Content = content
	task.Completed = completed
	s.tasks[id] = task

	return task, nil
}

func (s *TaskStore) DeleteTask(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.tasks[id]
	if !ok {
		return fmt.Errorf("task with id %d not found", id)
	}
	delete(s.tasks, id)
	return nil
}
