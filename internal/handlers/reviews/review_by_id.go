package reviews

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
)

func ReviewByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	reviewLookup := ps.ByName("id")

	reviewID, err := strconv.Atoi(reviewLookup)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
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
		r.id AS review_id,
		r.title AS review_title,
		r.content,
		r.total_likes,
		r.rating,
		r.spoiler,
		r.user_id,
		u.username,
		b.id AS book_id,
		b.title AS book_title,
		COALESCE(json_agg(json_build_object('id', a.id, 'name', a.name)) FILTER (WHERE a.id IS NOT NULL), '[]') AS authors
	FROM "Review" r
	JOIN "Document" b ON r.document_id = b.id
	LEFT JOIN "Author_Document" ad ON ad.document_id = b.id
	LEFT JOIN "Author" a ON ad.author_id = a.id
	LEFT JOIN "User" u ON r.user_id = u.id
	WHERE r.id = $1
	GROUP BY r.id, b.id, u.id
	`
	var review dbtypes.UniqueReview
	var authorsJSON string

	err = conn.QueryRow(context.Background(), query, reviewID).Scan(
		&review.Id,
		&review.Title,
		&review.Content,
		&review.Total_likes,
		&review.Rating,
		&review.Spoiler,
		&review.User_id,
		&review.Username,
		&review.Book.Id,
		&review.Book.Title,
		&authorsJSON,
	)
	if err == pgx.ErrNoRows {
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Error fetching review details:", err)
		http.Error(w, "Error fetching review details", http.StatusInternalServerError)
		return
	}

	urlImage := fmt.Sprintf("http://143.198.142.139:8080/cover/%d", review.Book.Id)
	review.Book.Cover_url = urlImage

	if err := json.Unmarshal([]byte(authorsJSON), &review.Book.Authors); err != nil {
		log.Println("Error parsing authors JSON:", err)
		http.Error(w, "Error processing review", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(review); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
