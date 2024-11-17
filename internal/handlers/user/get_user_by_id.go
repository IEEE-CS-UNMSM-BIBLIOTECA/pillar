package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
)

func GetUserById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userLookup := ps.ByName("user_id")

	username, hasToken := r.Context().Value("username").(string)
	var requesterID int
	if hasToken {
		requesterID = auth.GetIdFromUsername(username)
		if requesterID == 0 {
			http.Error(w, "That username does not exist", http.StatusBadRequest)
			return
		}
	}

	userID, err := strconv.Atoi(userLookup)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var user dbtypes.User
	query := `
	SELECT id, username, bio FROM "User" WHERE id = $1 
	`
	err = conn.QueryRow(context.Background(), query, userID).Scan(&user.Id, &user.Name, &user.Bio)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("User with ID %d not found\n", userID)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Println("Error scanning row:", err)
		http.Error(w, "Error processing users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
