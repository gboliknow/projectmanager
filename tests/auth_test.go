package test

import (
	"projectmanager/api"
	"testing"
)

func TestCreateJWT(t *testing.T) {
	secret := []byte("secret")
	userID := int64(50)

	t.Run("should create a valid JWT", func(t *testing.T) {
		tokenString, err := api.CreateJWT(secret, userID)

		if err != nil {
			t.Fatalf("expected no error , got %v", err)
		}

		if tokenString == "" {
			t.Errorf("expected token to not be empty")
		}
	})
}

func TestHashPassword(t *testing.T) {
	hash, err := api.HashPassword("password")

	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if len(hash) == 0 {
		t.Errorf("expected hash to no be empty")
	}

	if string(hash) == "password" {
		t.Errorf("expected hash to not be equal to password")
	}
}
