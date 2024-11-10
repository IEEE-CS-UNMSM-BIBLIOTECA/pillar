package handlers

import (
	"context"
	"log"
	"net/http"
	dbutils "pillar/internal/db/utils"
	"time"

	"github.com/go-json-experiment/json"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/julienschmidt/httprouter"
)

type SignUpRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	BirthDate   string `json:"birth_date"` // Expecting a date string
	Address     string `json:"address"`
	MobilePhone string `json:"mobile_phone"`
	RoleID      int32  `json:"role_id"`
	GenderID    int32  `json:"gender_id"`
}

func HndSignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req SignUpRequest
	if err := json.UnmarshalRead(r.Body, &req, json.DefaultOptionsV2()); err != nil {
		log.Println("Failed to parse request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" || req.Email == "" {
		w.WriteHeader(http.StatusBadRequest) // Bad request if essential fields are empty
		return
	}

	// Hash the password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Parse birth date
	birthDate, err := time.Parse("2006-01-02", req.BirthDate) // Expecting format YYYY-MM-DD
	if err != nil {
		log.Println("Invalid birth date format:", err)
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

	// Store the user in the database
	_, err = conn.Exec(context.Background(),
		"INSERT INTO \"User\" (username, email, bpassword, name, birth_date, address, mobile_phone, role_id, gender_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		req.Username, req.Email, hashedPassword, req.Name, birthDate, req.Address, req.MobilePhone, req.RoleID, req.GenderID)

	if err != nil {
		// Handle duplicate username or email error
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			log.Println("Conflict error on field:", pgerr.ConstraintName)
			w.WriteHeader(http.StatusConflict)
			json.MarshalWrite(w, map[string]string{"error": "Username or email already exists"}, json.DefaultOptionsV2())
			return
		}
		log.Println("Error inserting user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated) // Created
	json.MarshalWrite(w, map[string]string{"message": "User created successfully"}, json.DefaultOptionsV2())
}
