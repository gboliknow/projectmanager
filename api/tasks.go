package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
	r.HandleFunc("/tasks{id}", WithJWTAuth(s.handleGetTask, s.store)).Methods("GET")
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
	t, err := s.store.CreateTask(task)
	if err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utility.WriteJSON(w, http.StatusCreated, "Ok", t)

}

func (s *TasksService) handleGetTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	// if id == ""{
	// 	WriteJSON(w,http.StatusBadRequest, ErrorResponse{Error: "task not found"})
	// 	return
	// }

	t, err := s.store.GetTask(id)

	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "task not found", nil)
		return
	}

	utility.WriteJSON(w, http.StatusOK, "Ok", t)

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
