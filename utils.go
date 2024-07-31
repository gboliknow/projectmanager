package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)


func WriteJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    response := Response{
        StatusCode: statusCode,
        Message:    message,
        Data:       data,
    }
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
    }
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	headerParts := strings.Split(authHeader, " ")
	errTokenMissing := errors.New("missing or invalid token")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errTokenMissing
	}
	tokenAuth := headerParts[1]
	if tokenAuth != "" {
		return tokenAuth, nil
	}
	return "", errTokenMissing
}

