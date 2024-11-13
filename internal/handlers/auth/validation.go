package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
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

		// The token is passed as a "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Check if the header contains 'Bearer'
			http.Error(w, "Malformed authorization header", http.StatusUnauthorized)
			return
		}

		// Validate the JWT token
		token, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Optionally, check for specific claims (like user role or username) from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Extract the username from claims
		username := claims["username"].(string)
		log.Printf("Authenticated user: %s", username)

		// Set the username in context so it can be accessed by subsequent handlers
		ctx := context.WithValue(r.Context(), "username", username)
		r = r.WithContext(ctx)

		// Call the next handler in the chain
		next(w, r, ps)
	}
}

// Example of a protected endpoint that requires user authentication
func HndProtectedEndpoint(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve the username from the context
	username := r.Context().Value("username").(string)

	// You can now access the username for your business logic
	w.Write([]byte("Hello, " + username + ", you're authenticated!"))
}
