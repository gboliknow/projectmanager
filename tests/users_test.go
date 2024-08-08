package test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gorilla/mux"
// )

// func TestValidateUserPayload(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		user     *User
// 		expected error
// 	}{
// 		{
// 			name: "missing email",
// 			user: &User{FirstName: "John", LastName: "Doe", Password: "password123"},
// 			expected: errEmailRequired,
// 		},
// 		{
// 			name: "missing first name",
// 			user: &User{Email: "john@example.com", LastName: "Doe", Password: "password123"},
// 			expected: errFirstNameRequired,
// 		},
// 		{
// 			name: "missing last name",
// 			user: &User{Email: "john@example.com", FirstName: "John", Password: "password123"},
// 			expected: errLastNameRequired,
// 		},
// 		{
// 			name: "missing password",
// 			user: &User{Email: "john@example.com", FirstName: "John", LastName: "Doe"},
// 			expected: errPasswordRequired,
// 		},
// 		{
// 			name: "valid user",
// 			user: &User{Email: "john@example.com", FirstName: "John", LastName: "Doe", Password: "password123"},
// 			expected: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := validateUserPayload(tt.user)
// 			if err != tt.expected {
// 				t.Errorf("expected error %v, got %v", tt.expected, err)
// 			}
// 		})
// 	}
// }

// func TestCreateUser(t *testing.T) {
// 	// Create a new project
// 	ms := &MockStore{}
// 	service := NewUserService(ms)

// 	t.Run("should validate if the email is not empty", func(t *testing.T) {
// 		payload := &RegisterPayload{
// 			Email:     "",
// 			FirstName: "John",
// 			LastName:  "Doe",
// 			Password:  "password",
// 		}

// 		b, err := json.Marshal(payload)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(b))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		router := mux.NewRouter()

// 		router.HandleFunc("/users/register", service.handleUserRegister)

// 		router.ServeHTTP(rr, req)

// 		if rr.Code != http.StatusBadRequest {
// 			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
// 		}

// 		var response ErrorResponse
// 		err = json.NewDecoder(rr.Body).Decode(&response)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if response.Error != errEmailRequired.Error() {
// 			t.Errorf("expected error message %s, got %s", response.Error, errEmailRequired.Error())
// 		}
// 	})

// 	t.Run("should create a user", func(t *testing.T) {
// 		payload := &RegisterPayload{
// 			Email:     "joe@mail.com",
// 			FirstName: "John",
// 			LastName:  "Doe",
// 			Password:  "password",
// 		}

// 		b, err := json.Marshal(payload)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(b))
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		router := mux.NewRouter()

// 		router.HandleFunc("/users/register", service.handleUserRegister)

// 		router.ServeHTTP(rr, req)

// 		if rr.Code != http.StatusCreated {
// 			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
// 		}
// 	})
// }
