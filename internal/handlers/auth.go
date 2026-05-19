package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mike-testut/task-api/internal/httpjson"
	"github.com/mike-testut/task-api/internal/service"
)

type AuthHandlers struct {
	authService *service.AuthService
}

func NewAuthHandlers(as *service.AuthService) *AuthHandlers {
	return &AuthHandlers{authService: as}
}

type authInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input authInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		httpjson.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if input.Username == "" || input.Password == "" {
		httpjson.ErrorJSON(w, http.StatusBadRequest, "Username and password are required")
		return
	}
	user, err := h.authService.Register(input.Username, input.Password)
	if err != nil {
		httpjson.ErrorJSON(w, http.StatusConflict, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusCreated, user)
}

func (h *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input authInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		httpjson.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	tokenString, err := h.authService.Login(input.Username, input.Password)
	if err != nil {
		httpjson.ErrorJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	response := map[string]string{
		"token": tokenString,
	}
	httpjson.WriteJSON(w, http.StatusOK, response)
}
