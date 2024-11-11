package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	dbutils "pillar/internal/db/utils"
	"time"

	jsonexp "github.com/go-json-experiment/json"
	"github.com/golang-jwt/jwt/v5"
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

func HndLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var creds Credentials
	if err := jsonexp.UnmarshalRead(r.Body, &creds, jsonexp.DefaultOptionsV2()); err != nil {
		log.Println("Failed to parse request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Acquire a database connection
	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	// Prepare and execute the query to find the user
	var storedPassword string
	err = conn.QueryRow(context.Background(),
		"SELECT bpassword FROM \"User\" WHERE username = $1", creds.Username).Scan(&storedPassword)

	var tag []string
	docs, errores := conn.Query(context.Background(),
		"SELECT title FROM \"Document\"")
	if errores != nil {
		log.Println("Query error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for docs.Next() {
		var tagID string
		log.Printf("ALGO:")
		if err := docs.Scan(&tagID); err != nil {
			log.Printf("Error scanning row: %v", err)
			return
		}

		tag = append(tag, tagID)
	}
	log.Printf("AYUDA: %v", tag)
	// Check for errors
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Println("Query error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Verify the password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized) // Password incorrect
		return
	}

	// Generate JWT if credentials are valid
	token, err := GenerateJWT(creds.Username)
	if err != nil {
		log.Println("Failed to generate token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the token back to the client
	w.Header().Set("Content-Type", "application/json")
	jsonexp.MarshalWrite(w, map[string]string{"token": token}, jsonexp.DefaultOptionsV2())
}
