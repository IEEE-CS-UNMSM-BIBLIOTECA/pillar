package normal_user

import (
	"context"
	"log"
	"net/http"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"

	jsonexp "github.com/go-json-experiment/json"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

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
	token, err := auth.GenerateJWT(creds.Username)
	if err != nil {
		log.Println("Failed to generate token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the token back to the client
	w.Header().Set("Content-Type", "application/json")
	jsonexp.MarshalWrite(w, map[string]string{"token": token}, jsonexp.DefaultOptionsV2())
}
