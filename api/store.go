package api

import (
	"database/sql"
	"fmt"
	"projectmanager/internal/types"
)

type Store interface {
	// Users
	GetUserByID(id string) (*types.User, error)
	GetUserByEmail(email string) (*types.User, error)
	CreateUser(u *types.User) (*types.User, error)
	//Tasks
	CreateTask(t *types.Task) (*types.Task, error)
	GetTask(id string) (*types.Task, error)

	//Projects
	CreateProject(p *types.Project) error
	GetProject(id string) (*types.Project, error)
	DeleteProject(id string) error
	GetProjectByName(name string) (bool, error)
	GetAllProjects() ([]*types.Project, error)
	DeleteAllProjects() error
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(u *types.User) (*types.User, error) {
	rows, err := s.db.Exec("INSERT INTO users (email, firstName, lastName, password) VALUES (?, ?, ?, ?)", u.Email, u.FirstName, u.LastName, u.Password)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = id
	return u, nil
}

func (s *Storage) CreateTask(t *types.Task) (*types.Task, error) {

	if t.Status == "" {
		t.Status = "TODO" // default value
	}
	rows, err := s.db.Exec("INSERT INTO tasks (name, status, projectId, AssignedToID) VALUES (?, ?, ?, ?)", t.Name, t.Status, t.ProjectID, t.AssignedToID)

	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()

	if err != nil {
		return nil, err
	}

	t.ID = id
	return t, nil
}

func (s *Storage) GetTask(id string) (*types.Task, error) {
	var t types.Task
	err := s.db.QueryRow("SELECT id, name, status, projectId, AssignedToID, createdAt FROM tasks WHERE id = ?", id).Scan(&t.ID, &t.Name, &t.Status, &t.ProjectID, &t.AssignedToID, &t.CreatedAt)
	return &t, err
}

func (s *Storage) GetUserByID(id string) (*types.User, error) {
	var u types.User
	err := s.db.QueryRow("SELECT id, email, firstName, lastName, createdAt FROM users WHERE id = ?", id).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt)
	return &u, err
}

func (s *Storage) GetUserByEmail(email string) (*types.User, error) {
	var u types.User
	err := s.db.QueryRow("SELECT id, email, firstName, lastName, createdAt, password FROM users WHERE email = ?", email).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (s *Storage) CreateProject(p *types.Project) error {
	result, err := s.db.Exec("INSERT INTO projects (name) VALUES (?)", p.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	p.ID = id
	return nil
}

func (s *Storage) GetProject(id string) (*types.Project, error) {
	var p types.Project
	err := s.db.QueryRow("SELECT id, name, createdAt FROM projects WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.CreatedAt)
	return &p, err
}

func (s *Storage) GetAllProjects() ([]*types.Project, error) {
	rows, err := s.db.Query("SELECT id, name, createdAt FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*types.Project
	for rows.Next() {
		var p types.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *Storage) DeleteProject(id string) error {
	_, err := s.db.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteAllProjects() error {
	_, err := s.db.Exec("DELETE FROM projects")
	return err
}

func (s *Storage) GetProjectByName(name string) (bool, error) {
	var p types.Project
	err := s.db.QueryRow("SELECT id, name, createdAt FROM projects WHERE name = ?", name).Scan(&p.ID, &p.Name, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return false, nil // No project found with the given name
	}
	if err != nil {
		return false, err // Other errors (e.g., query issues)
	}
	return true, nil // Project found
}
