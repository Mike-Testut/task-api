package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mike-testut/task-api/internal/service"
)

type TaskHandlers struct {
	service *service.TaskService
}

func NewTaskHandlers(s *service.TaskService) *TaskHandlers {
	return &TaskHandlers{service: s}
}

func (h *TaskHandlers) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.service.CreateTask(input.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandlers) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTask((id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandlers) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks,_ := h.service.ListTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandlers) UpdateTaskHandler(w http.ResponseWriter, r *http.Request){
	idStr := r.PathValue("id")
	id,err := strconv.Atoi(idStr)
	if err!=nil{
		http.Error(w,"Invalid task ID", http.StatusBadRequest)
		return
	}
	var input struct{
		Content string `json:"content"`
		Completed bool `json:"completed"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err !=nil{
		http.Error(w,err.Error(), http.StatusBadRequest)
		return
	}
	task, err := h.service.UpdateTask(id, input.Content, input.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request){
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteTask(id)
	if err != nil{
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}