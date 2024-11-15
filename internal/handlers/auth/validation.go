package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	dbutils "pillar/internal/db/utils"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("teemo")

func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 1 día de expiración

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metodo de firma inesperado")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// Middleware to validate JWT token for each protected endpoint
func TokenValidationMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Extract the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Check if the header contains 'Bearer'
			http.Error(w, "Malformed authorization header", http.StatusUnauthorized)
			return
		}

		token, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		username := claims["username"].(string)
		log.Printf("Authenticated user: %s", username)

		ctx := context.WithValue(r.Context(), "username", username)
		r = r.WithContext(ctx)

		next(w, r, ps)
	}
}

func HndProtectedEndpoint(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve the username from the context
	username := r.Context().Value("username").(string)

	w.Write([]byte("Hello, " + username + ", you're authenticated!"))
}

func GetIdFromUsername(username string) int {
	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		return 0
	}
	defer conn.Release()

	var userID int
	query := `SELECT id FROM "User" WHERE username = $1`
	err = conn.QueryRow(context.Background(), query, username).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("No user found with the given username")
		} else {
			log.Println("Error executing query:", err)
		}
		return 0
	}

	return userID
}
