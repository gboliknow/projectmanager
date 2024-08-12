package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"projectmanager/internal/utility"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := utility.GetTokenFromRequest(r)
		if err != nil {
			errorHandler(w, "missing or invalid token")
			return
		}

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to authenticate token with err : %v", err)
			errorHandler(w, "permission denied")
			return
		}

		if !token.Valid {
			log.Printf("failed to authenticate token because it invalid")
			errorHandler(w, "permission denied")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(string)

		_, err = store.GetUserByID(userID)

		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			errorHandler(w, "permission denied")
			return
		}
		handlerFunc(w, r)
	}
}

func errorHandler(w http.ResponseWriter, errorString string) {
	utility.WriteJSON(w, http.StatusUnauthorized, errorString, nil)
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func CreateJWT(secret []byte, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
