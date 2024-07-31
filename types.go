package main

import "time"

type ErrorResponse struct {
	Error string `JSON:"error"`
}

type Task struct {
	ID           int64     `JSON:"id"`
	Name         string    `JSON:"name"`
	Status       string    `JSON:"status"`
	ProjectID    int64     `JSON:"projectID"`
	AssignedToID int64     `JSON:"assignedTo"`
	CreatedAt    time.Time `JSON:"createdAt"`
}

type User struct {
	ID        int64     `JSON:"id"`
	FirstName string    `JSON:"firstName"`
	Email     string    `JSON:"email"`
	LastName  string    `JSON:"lastName"`
	Password  string    `JSON:"password"`
	CreatedAt time.Time `JSON:"createdAt"`
}


type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type Project struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateProjectPayload struct {
	Name string `json:"name"`
}

type Response struct {
    StatusCode int         `json:"statusCode"`
    Message    string      `json:"message"`
    Data       interface{} `json:"data,omitempty"` // Data is omitted if nil or empty
}