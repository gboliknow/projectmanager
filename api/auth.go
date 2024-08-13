package api

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
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
			log.Printf("Failed to authenticate token: %v", err)
			errorHandler(w, "permission denied")
			return
		}

		if !token.Valid {
			log.Printf("Token is invalid")
			errorHandler(w, "permission denied")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("Invalid token claims")
			errorHandler(w, "permission denied")
			return
		}

		userIDStr, ok := claims["userID"].(string)
		if !ok {
			log.Printf("UserID not found or invalid in token")
			errorHandler(w, "permission denied")
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Printf("Failed to parse userID from token: %v", err)
			errorHandler(w, "permission denied")
			return
		}

		_, err = store.GetUserByID(userID)
		if err != nil {
			log.Printf("Failed to get user by ID: %v", err)
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
	// secret := os.Getenv("JWT_SECRET")
	secret := "c757b8c7cacc1d63b3d37a5688eaef1809687c1d3a4330192c7fbfe93a8dbeb5"
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

func getUserIDFromToken(tokenString string, secret []byte) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	userIDStr, ok := claims["userID"].(string)
	if !ok {
		return 0, fmt.Errorf("userID not found in token")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid userID format")
	}

	return userID, nil
}

func sendPasswordResetEmail(email, token string) error {
	from := "gboliknow@gmail.com"
	password := "*********"
	fmt.Println(token)
	to := email
	subject := "Password Reset Request"
	body := fmt.Sprintf("To reset your password, use the following token: %s", token)

	msg := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))

	auth := smtp.PlainAuth("", from, password, "smtp.example.com")
	err := smtp.SendMail("smtp.example.com:587", auth, from, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}


func validatePassword(password string) error {
    if len(password) < 8 {
        return fmt.Errorf("password must be at least 8 characters long")
    }
    // Add more checks for password strength as needed
    return nil
}