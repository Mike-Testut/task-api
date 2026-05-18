package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mike-testut/task-api/internal/httpjson"
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
		httpjson.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.service.CreateTask(input.Content)
	if err != nil {
		if errors.Is(err, service.ErrContentRequired) || errors.Is(err, service.ErrContentTooLong) {
			httpjson.ErrorJSON(w, http.StatusBadRequest, err.Error())
		} else {
			httpjson.ErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, task)
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
		if errors.Is(err, service.ErrTaskNotFound) {
			httpjson.ErrorJSON(w, http.StatusNotFound, err.Error())
		} else {
			httpjson.ErrorJSON(w, http.StatusInternalServerError, "Internal service error")
		}
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, task)
}

func (h *TaskHandlers) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limitStr := query.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 25
	}

	offsetStr := query.Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	tasks, err := h.service.ListTasks(limit, offset)
	if err != nil {
		httpjson.ErrorJSON(w, http.StatusInternalServerError, "Interna server error")
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandlers) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		httpjson.ErrorJSON(w, http.StatusBadRequest, "invalid task ID")
		return
	}
	var input struct {
		Content   string `json:"content"`
		Completed bool   `json:"completed"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		httpjson.ErrorJSON(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	task, err := h.service.UpdateTask(id, input.Content, input.Completed)
	if err != nil {
		httpjson.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, task)
}

func (h *TaskHandlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
