package api

import (
	"database/sql"
	"fmt"
	"projectmanager/internal/types"
	"time"
)

type Store interface {
	// Users
	GetUserByID(id int64) (*types.User, error)
	GetUserByEmail(email string) (*types.User, error)
	CreateUser(u *types.User) (*types.User, error)
	UpdateUserProfile(userID int64, updateRequest *types.UserUpdateRequest) (*types.User, error)
	UpdatePassword(userID int64, newPassword string) error
	ValidateResetToken(token string) (int64, error)
	InvalidateResetToken(token string) error
	RequestPasswordReset(email, resetToken string) error

	//Tasks
	CreateTask(t *types.Task) (*types.Task, error)
	GetTask(id string) (*types.Task, error)
	TaskExists(t *types.Task) (bool, error)
	GetMyTasks(userID int64, status string) ([]types.Task, error)

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

func (s *Storage) TaskExists(t *types.Task) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM tasks WHERE name = ? AND projectID = ?"
	err := s.db.QueryRow(query, t.Name, t.ProjectID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (s *Storage) GetTask(id string) (*types.Task, error) {
	var t types.Task
	err := s.db.QueryRow("SELECT id, name, status, projectId, AssignedToID, createdAt FROM tasks WHERE id = ?", id).Scan(&t.ID, &t.Name, &t.Status, &t.ProjectID, &t.AssignedToID, &t.CreatedAt)
	return &t, err
}

func (s *Storage) GetMyTasks(userID int64, status string) ([]types.Task, error) {
	var tasks []types.Task
	var rows *sql.Rows
	var err error
	if status == "" {
		rows, err = s.db.Query("SELECT id, name, status, projectId, AssignedToID, createdAt FROM tasks WHERE AssignedToID = ?", userID)
	} else {
		rows, err = s.db.Query("SELECT id, name, status, projectId, AssignedToID, createdAt FROM tasks WHERE AssignedToID = ? AND status = ?", userID, status)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t types.Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Status, &t.ProjectID, &t.AssignedToID, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Storage) GetUserByID(id int64) (*types.User, error) {
	var u types.User
	err := s.db.QueryRow("SELECT id, email, firstName, lastName, phone, address,password, createdAt FROM users WHERE id = ?", id).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.Address, &u.Password, &u.CreatedAt)
	return &u, err
}

func (s *Storage) GetUserByEmail(email string) (*types.User, error) {
	var u types.User
	err := s.db.QueryRow("SELECT id, email, firstName, lastName, phone, address,password, createdAt FROM users WHERE email = ?", email).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.Address, &u.Password, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (s *Storage) UpdateUserProfile(userID int64, updateRequest *types.UserUpdateRequest) (*types.User, error) {
	query := "UPDATE users SET"
	args := []interface{}{}
	argCount := 1

	if updateRequest.FirstName != nil {
		query += " firstName = ?"
		args = append(args, *updateRequest.FirstName)
		argCount++
	}
	if updateRequest.LastName != nil {
		if argCount > 1 {
			query += ","
		}
		query += " lastName = ?"
		args = append(args, *updateRequest.LastName)
		argCount++
	}
	if updateRequest.Email != nil {
		if argCount > 1 {
			query += ","
		}
		query += " email = ?"
		args = append(args, *updateRequest.Email)
		argCount++
	}
	// Add more fields similarly
	if updateRequest.Phone != nil {
		if argCount > 1 {
			query += ","
		}
		query += " phone = ?"
		args = append(args, *updateRequest.Phone)
		argCount++
	}
	if updateRequest.Address != nil {
		if argCount > 1 {
			query += ","
		}
		query += " address = ?"
		args = append(args, *updateRequest.Address)
		argCount++
	}

	query += " WHERE id = ?"
	args = append(args, userID)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// Fetch the updated user details
	var updatedUser types.User
	err = s.db.QueryRow("SELECT id, email, firstName, lastName, phone, address, createdAt FROM users WHERE id = ?", userID).Scan(&updatedUser.ID, &updatedUser.Email, &updatedUser.FirstName, &updatedUser.LastName, &updatedUser.Phone, &updatedUser.Address, &updatedUser.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
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

func (s *Storage) RequestPasswordReset(email, resetToken string) error {
	// Check if the email exists
	var userID int64
	err := s.db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("email not found")
		}
		return err
	}

	// Save the reset token and expiration to the database
	expiration := time.Now().Add(1 * time.Hour) // Token valid for 1 hour
	_, err = s.db.Exec(`
		INSERT INTO password_reset_tokens (user_id, token, expiration)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE token = VALUES(token), expiration = VALUES(expiration)
	`, userID, resetToken, expiration)

	return err
}

func (s *Storage) UpdatePassword(userID int64, newPassword string) error {
	_, err := s.db.Exec(`
		UPDATE users
		SET password = ?
		WHERE id = ?
	`, newPassword, userID)
	return err
}

func (s *Storage) ValidateResetToken(token string) (int64, error) {
	var userID int64
	var expiration time.Time

	err := s.db.QueryRow(`
		SELECT user_id, expiration
		FROM password_reset_tokens
		WHERE token = ?
	`, token).Scan(&userID, &expiration)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("invalid token")
		}
		return 0, err
	}

	if time.Now().After(expiration) {
		return 0, fmt.Errorf("token expired")
	}

	return userID, nil
}

func (s *Storage) InvalidateResetToken(token string) error {
	_, err := s.db.Exec("DELETE FROM password_reset_tokens WHERE token = ?", token)
	return err
}
