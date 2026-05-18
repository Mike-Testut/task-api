package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mike-testut/task-api/internal/service"
	"github.com/mike-testut/task-api/internal/store"
)

func TestListTasksHandler(t *testing.T){
	taskStore := store.NewTaskStore()
	taskService := service.NewTaskService(taskStore)
	taskHandler := NewTaskHandlers(taskService)

	taskService.CreateTask("Test Task 1")

	req:=httptest.NewRequest("Get", "/tasks", nil)
	rr:= httptest.NewRecorder()

	taskHandler.ListTasksHandler(rr,req)

	if rr.Code != http.StatusOK{
		t.Errorf("expected status %d; got %d", http.StatusOK, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Test Task 1"){
		t.Errorf("response body does not contain expected task")
	}

	if rr.Header().Get("Content-Type") != "application/json"{
		t.Errorf("expected content type json; got %s", rr.Header().Get("Content-Type"))
	}
}

