package store

import (
	"fmt"
	"sync"
	"github.com/mike-testut/task-api/internal/models"
)

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

func (s *TaskStore) CreateTask(content string) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := models.Task{
		ID:        s.nextID,
		Content:   content,
		Completed: false,
	}

	s.tasks[task.ID] = task
	s.nextID++

	return task
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

func(s *TaskStore) ListTasks()[]models.Task{
	s.mu.Lock()
	defer s.mu.Unlock()

	allTasks := make([]models.Task, 0, len(s.tasks))

	for _,task := range s.tasks{
		allTasks = append(allTasks, task)
	}
	return allTasks
}
