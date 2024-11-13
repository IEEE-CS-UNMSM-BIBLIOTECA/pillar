package admin

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

func HndLoginAdmin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var storedPassword string
	var roleID int
	err = conn.QueryRow(context.Background(),
		"SELECT bpassword, role_id FROM \"User\" WHERE username = $1", creds.Username).Scan(&storedPassword, &roleID)

	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Println("Query error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if the role_id is 1 (admin)
	if roleID != 1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	token, err := auth.GenerateJWT(creds.Username)
	if err != nil {
		log.Println("Failed to generate token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonexp.MarshalWrite(w, map[string]string{"token": token}, jsonexp.DefaultOptionsV2())
}
