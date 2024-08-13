package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"projectmanager/internal/config"
	"projectmanager/internal/types"
	"projectmanager/internal/utility"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var errEmailRequired = errors.New("email is required")
var errFirstNameRequired = errors.New("first name is required")
var errLastNameRequired = errors.New("last name is required")
var errPasswordRequired = errors.New("password is required")

type UserService struct {
	store Store
}

func NewUserService(s Store) *UserService {
	return &UserService{store: s}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.handleUserRegister).Methods("POST")
	r.HandleFunc("/users/login", s.handleUserLogin).Methods("POST")
	r.HandleFunc("/users/me", WithJWTAuth(s.handleUpdateUserProfile, s.store)).Methods("PUT")
	r.HandleFunc("/users/me", s.handleGetUserInfo).Methods("GET")
	r.HandleFunc("/users/reset-password", s.handlePasswordResetRequest).Methods("POST")
	r.HandleFunc("/users/reset-password/confirm", s.handleResetPassword).Methods("POST")
	// r.HandleFunc("/users/logout", s.handleUpdateUserProfile).Methods("POST")
}

func (s *UserService) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload *types.User
	err = json.Unmarshal(body, &payload)
	if err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if err := validateUserPayload(payload); err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	hashedPassword, err := HashPassword(payload.Password)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error creating user", nil)
		return
	}
	payload.Password = hashedPassword

	u, err := s.store.CreateUser(payload)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error creating user", nil)
		return
	}

	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error creating user", nil)
		return
	}

	utility.WriteJSON(w, http.StatusCreated, "Successful", token)
}

func (s *UserService) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var payload *types.User
	err = json.Unmarshal(body, &payload)
	if err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if err := validateLoginUserPayload(payload); err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)

		return
	}

	user, err := s.store.GetUserByEmail(payload.Email)
	if err != nil {
		if err.Error() == "user not found" {
			utility.WriteJSON(w, http.StatusUnauthorized, "User not found", nil)

		} else {
			utility.WriteJSON(w, http.StatusInternalServerError, "Error retrieving user", nil)
		}
		return
	}
	fmt.Printf("the user payload %v", user)
	if !CheckPasswordHash(payload.Password, user.Password) {
		utility.WriteJSON(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	token, err := createAndSetAuthCookie(user.ID, w)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error logging in", nil)
		return
	}

	responseData := struct {
		Token string              `json:"token"`
		User  *types.UserResponse `json:"user"`
	}{
		Token: token,
		User: &types.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			Address:   user.Address,
			Phone:     user.Phone,
		},
	}
	utility.WriteJSON(w, http.StatusOK, "Successful", responseData)
}

func (s *UserService) handleGetUserInfo(w http.ResponseWriter, r *http.Request) {
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
		utility.WriteJSON(w, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	user, err := s.store.GetUserByID(userID)
	responseData := types.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Address:   user.Address,
		Phone:     user.Phone,
	}

	if err != nil {
		fmt.Println(err.Error())
		utility.WriteJSON(w, http.StatusInternalServerError, "Error Fetching User", nil)
		return
	}

	utility.WriteJSON(w, http.StatusOK, "User Fetched Successfully", responseData)
}

func (s *UserService) handleUpdateUserProfile(w http.ResponseWriter, r *http.Request) {
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
		utility.WriteJSON(w, http.StatusUnauthorized, "Invalid token", nil)
		return
	}

	var payload types.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	user, err := s.store.UpdateUserProfile(userID, &payload)
	responseData := types.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Address:   user.Address,
		Phone:     user.Phone,
	}
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error updating profile", nil)
		return
	}

	utility.WriteJSON(w, http.StatusOK, "Profile updated", responseData)
}

func (s *UserService) handlePasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	var payload types.PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	// Generate a reset token
	resetToken := utility.GenerateResetToken()

	// Store the token and expiration
	err := s.store.RequestPasswordReset(payload.Email, resetToken)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error processing request", nil)
		return
	}

	// Send email with the reset token
	err = sendPasswordResetEmail(payload.Email, resetToken)
	if err != nil {
		fmt.Println(err.Error())
		utility.WriteJSON(w, http.StatusInternalServerError, "Error sending email", nil)
		return
	}

	utility.WriteJSON(w, http.StatusOK, "Password reset email sent", nil)
}

func (s *UserService) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	var payload types.PasswordResetPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	// Validate that required fields are not empty
	if payload.ResetToken == "" || payload.NewPassword == "" {
		utility.WriteJSON(w, http.StatusBadRequest, "Reset token and new password are required", nil)
		return
	}

	// Validate the reset token and get user ID
	userID, err := s.store.ValidateResetToken(payload.ResetToken)
	if err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, "Invalid or expired reset token", nil)
		return
	}

	// Validate the new password
	if err := validatePassword(payload.NewPassword); err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, "Weak password", nil)
		return
	}

	// Hash the new password
	hashedPassword, err := HashPassword(payload.NewPassword)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error hashing password", nil)
		return
	}

	// Update the user's password
	err = s.store.UpdatePassword(userID, hashedPassword)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error updating password", nil)
		return
	}

	// Invalidate the reset token
	err = s.store.InvalidateResetToken(payload.ResetToken)
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, "Error invalidating reset token", nil)
		return
	}

	utility.WriteJSON(w, http.StatusOK, "Password updated successfully", nil)
}

// func (s *UserService) handleLogout(w http.ResponseWriter, r *http.Request) {
// 	// Handle token invalidation or session management here
// 	utility.WriteJSON(w, http.StatusOK, "Logout successful", nil)
// }

func validateUserPayload(user *types.User) error {
	if user.Email == "" {
		return errEmailRequired
	}

	if user.FirstName == "" {
		return errFirstNameRequired
	}

	if user.LastName == "" {
		return errLastNameRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}

func validateLoginUserPayload(user *types.User) error {
	if user.Email == "" {
		return errEmailRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}

func createAndSetAuthCookie(userID int64, w http.ResponseWriter) (string, error) {
	secret := []byte(config.Envs.JWTSecret)
	token, err := CreateJWT(secret, userID)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
