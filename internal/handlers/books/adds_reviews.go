package books

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"

	"github.com/julienschmidt/httprouter"
)

func AddReviews(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Context().Value("username").(string)

	user_id := auth.GetIdFromUsername(username)
	if user_id == 0 {
		http.Error(w, "That username does not exist", http.StatusBadRequest)
		return
	}

	var reviewReq dbtypes.ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&reviewReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
		`SELECT add_review($1, $2, $3, $4, $5, $6)`,
		user_id,
		reviewReq.DocumentID,
		reviewReq.Title,
		reviewReq.Content,
		reviewReq.Rating,
		reviewReq.Spoiler,
	)

	if err != nil {
		log.Println("Failed to execute add_review function:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Review added successfully"}`))
}
