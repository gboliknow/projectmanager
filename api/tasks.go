package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"projectmanager/internal/config"
	"projectmanager/internal/types"
	"projectmanager/internal/utility"

	"github.com/gorilla/mux"
)

var errNameRequired = errors.New("name is required")
var errProjectIDRequired = errors.New("project id is required")
var errUserIDRequired = errors.New("user id is required")

type TasksService struct {
	store Store
}

func NewTasksService(s Store) *TasksService {
	return &TasksService{store: s}
}

func (s *TasksService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", WithJWTAuth(s.handleCreateTask, s.store)).Methods("POST")
	r.HandleFunc("/tasks/{id}", WithJWTAuth(s.handleGetTask, s.store)).Methods("GET")
	r.HandleFunc("/mytasks", WithJWTAuth(s.handleGetMyTasks, s.store)).Methods("POST")
}

func (s *TasksService) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var task *types.Task

	err = json.Unmarshal(body, &task)

	if err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, "Invalid request payload!", nil)
		return
	}

	if err := validateTaskPayload(task); err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	exists, err := s.store.TaskExists(task)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error checking task existence", nil)
		return
	}
	if exists {
		utility.WriteJSON(w, http.StatusConflict, "Task already exists", nil)
		return
	}

	t, err := s.store.CreateTask(task)
	if err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utility.WriteJSON(w, http.StatusCreated, "Task created successfully", t)

}

func (s *TasksService) handleGetTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	t, err := s.store.GetTask(id)

	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "task not found", nil)
		return
	}

	utility.WriteJSON(w, http.StatusOK, "Ok", t)

}

func (s *TasksService) handleGetMyTasks(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var requestPayload struct {
		Status string `json:"status,omitempty"`
	}
	err = json.Unmarshal(body, &requestPayload)
	if err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	tokenString, err := utility.GetTokenFromRequest(r)
	if err != nil {
		errorHandler(w, "missing or invalid token")
		return
	}
	if tokenString == "" {
		utility.WriteJSON(w, http.StatusUnauthorized, "Missing token", nil)
		return
	}
	secret := []byte(config.Envs.JWTSecret)
	userID, err := getUserIDFromToken(tokenString, secret)

	if err != nil {
		utility.WriteJSON(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	tasks, err := s.store.GetMyTasks(userID, requestPayload.Status)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error retrieving tasks", nil)
		return
	}
	if tasks == nil {
		tasks = []types.Task{}
	}
	utility.WriteJSON(w, http.StatusOK, "Tasks retrieved successfully", tasks)
}


func validateTaskPayload(task *types.Task) error {
	if task.Name == "" {
		return errNameRequired
	}

	if task.ProjectID == 0 {
		return errProjectIDRequired
	}

	if task.AssignedToID == 0 {
		return errUserIDRequired
	}

	return nil
}
