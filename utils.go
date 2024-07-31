package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
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

