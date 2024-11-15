package normal_user

import (
	"context"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"

	"time"

	"github.com/go-json-experiment/json"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/julienschmidt/httprouter"
)

func HndSignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req dbtypes.SignUpRequest
	if err := json.UnmarshalRead(r.Body, &req, json.DefaultOptionsV2()); err != nil {
		log.Println("Failed to parse request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" || req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// The best way to security WAOS
	req.RoleID = 2

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate) // Expecting format YYYY-MM-DD
	if err != nil {
		log.Println("Invalid birth date format:", err)
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

	_, err = conn.Exec(context.Background(),
		"INSERT INTO \"User\" (username, email, bpassword, name, birth_date, address, mobile_phone, role_id, gender_id, bio, profile_picture_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		req.Username, req.Email, hashedPassword, req.Name, birthDate, req.Address, req.MobilePhone, req.RoleID, req.GenderID, req.Bio, req.Profile_picture_url)

	if err != nil {
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

	w.WriteHeader(http.StatusCreated)
	json.MarshalWrite(w, map[string]string{"message": "User created successfully"}, json.DefaultOptionsV2())
}
