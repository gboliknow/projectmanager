package api

import (
	"REST_API_WITH_GO/internal/types"
	"REST_API_WITH_GO/internal/utility"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type ProjectService struct {
	store Store
}

func NewProjectService(s Store) *ProjectService {
	return &ProjectService{store: s}
}

func (s *ProjectService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/projects", WithJWTAuth(s.handleCreateProject, s.store)).Methods("POST")
	r.HandleFunc("/projects/{id}", WithJWTAuth(s.handleGetProject, s.store)).Methods("GET")
	r.HandleFunc("/projects", WithJWTAuth(s.handleGetAllProjects, s.store)).Methods("GET")
	r.HandleFunc("/projects/{id}", WithJWTAuth(s.handleDeleteProject, s.store)).Methods("DELETE")
	r.HandleFunc("/projects", WithJWTAuth(s.handleDeleteAllProjects, s.store)).Methods("DELETE")
}

func (s *ProjectService) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	var project  *types.Project
	err = json.Unmarshal(body, &project)

	if err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if project.Name == "" {
		utility.WriteJSON(w, http.StatusBadRequest, "Name is required", nil)
		return
	}
	prjExist, err := s.store.GetProjectByName(project.Name)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error checking project existence", nil)
		return
	}

	if prjExist {
		utility.WriteJSON(w, http.StatusConflict, "Project with this name already exists", nil)
		return
	}
	err = s.store.CreateProject(project)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error creating project", nil)
		return
	}

	utility.WriteJSON(w, http.StatusCreated, "Ok", project)
}

func (s *ProjectService) handleGetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	project, err := s.store.GetProject(id)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error getting project", nil)
		return
	}
	utility.WriteJSON(w, http.StatusOK, "Ok", project)
}

func (s *ProjectService) handleGetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := s.store.GetAllProjects()
	if err != nil {
		utility.WriteJSON(w, http.StatusOK, "No projects", nil)
		return
	}
	if len(projects) == 0 {
		utility.WriteJSON(w, http.StatusOK, "No projects", []interface{}{})
		return
	}

	utility.WriteJSON(w, http.StatusOK, "Projects retrieved successfully", projects)
}
func (s *ProjectService) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := s.store.DeleteProject(id)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error deleting project", nil)
		return
	}

	utility.WriteJSON(w, http.StatusNoContent, "Projects Deleted successfully", nil)
}

func (s *ProjectService) handleDeleteAllProjects(w http.ResponseWriter, r *http.Request) {
	err := s.store.DeleteAllProjects()
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	 utility. WriteJSON(w, http.StatusOK, "All projects deleted successfully", nil)
}
