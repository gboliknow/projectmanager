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
