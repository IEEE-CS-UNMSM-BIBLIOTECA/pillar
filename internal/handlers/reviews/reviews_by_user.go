package reviews

import (
	"context"
	"encoding/json"
	"fmt"
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

	queryUserExists := `SELECT EXISTS(SELECT 1 FROM "User" WHERE id = $1)`
	var userExists bool
	err = conn.QueryRow(context.Background(), queryUserExists, userID).Scan(&userExists)
	if err != nil {
		log.Println("Error checking if user exists:", err)
		http.Error(w, "Error verifying user existence", http.StatusInternalServerError)
		return
	}

	if !userExists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	query := `
	SELECT
		r.id AS review_id,
		r.title AS review_title,
		r.content,
		r.total_likes,
		r.rating,
		r.spoiler,
		CASE WHEN $1 = r.user_id THEN true ELSE false END AS own,
		b.id AS book_id,
		b.title AS book_title,
		COALESCE(json_agg(json_build_object('id', a.id, 'name', a.name)) FILTER (WHERE a.id IS NOT NULL), '[]') AS authors
	FROM "Review" r
	JOIN "Document" b ON r.document_id = b.id
	LEFT JOIN "Author_Document" ad ON ad.document_id = b.id
	LEFT JOIN "Author" a ON ad.author_id = a.id
	WHERE r.user_id = $2
	GROUP BY r.id, b.id
	`
	rows, err := conn.Query(context.Background(), query, requesterID, userID)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching reviews", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reviews []dbtypes.ReviewUser
	for rows.Next() {
		var review dbtypes.ReviewUser
		var authorsJSON string

		err = rows.Scan(
			&review.Id,
			&review.Title,
			&review.Content,
			&review.Total_likes,
			&review.Rating,
			&review.Spoiler,
			&review.Own,
			&review.Book.Id,
			&review.Book.Title,
			&authorsJSON,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing reviews", http.StatusInternalServerError)
			return
		}

		urlImage := fmt.Sprintf("http://143.198.142.139:8080/cover/%d", review.Book.Id)
		review.Book.Cover_url = urlImage

		if err := json.Unmarshal([]byte(authorsJSON), &review.Book.Authors); err != nil {
			log.Println("Error parsing authors JSON:", err)
			http.Error(w, "Error processing reviews", http.StatusInternalServerError)
			return
		}

		reviews = append(reviews, review)
	}

	if len(reviews) == 0 {
		http.Error(w, "No reviews found", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
