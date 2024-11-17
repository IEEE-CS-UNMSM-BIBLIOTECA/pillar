package reviews

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func GetReviewsByUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	query := `
	SELECT 
		r.id,
		r.title,
		r.content,
		r.total_likes,
		r.rating,
		r.spoiler,
		CASE WHEN $1 = r.user_id THEN true ELSE false END AS own,
	FROM "Review" r
	JOIN "User" u ON r.user_id = u.id
	WHERE r.user_id = $2
	`
	rows, err := conn.Query(context.Background(), query, requesterID, userID)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching reviews", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reviews []dbtypes.Review
	for rows.Next() {
		var review dbtypes.Review
		err = rows.Scan(
			&review.Id,
			&review.Title,
			&review.Content,
			&review.Total_likes,
			&review.Rating,
			&review.Spoiler,
			&review.Own,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing reviews", http.StatusInternalServerError)
			return
		}
		reviews = append(reviews, review)
	}

	if len(reviews) == 0 {
		http.Error(w, "No reviews found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
