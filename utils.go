package main

import (
	"encoding/json"
	"net/http"
)



func WriteJSON(w http.ResponseWriter, status int, v any){
	w.Header().Set("Content-type" , "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}


func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}