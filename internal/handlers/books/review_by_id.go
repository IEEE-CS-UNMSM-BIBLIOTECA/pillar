package books

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type response struct {
	Review dbtypes.Review         `json:"review"`
	User   dbtypes.UserFromReview `json:"user"`
}

// Updated UserFromReview struct
type UserFromReview struct {
	Id                int32  `db:"id"`
	Name              string `db:"name"`
	ProfilePictureUrl string `db:"profile_picture_url"`
}

func fetchBookReviews(conn *pgxpool.Conn, bookID, page, pageSize int) ([]response, error) {
	query := `SELECT * FROM get_book_reviews($1, $2, $3)`
	rows, err := conn.Query(context.Background(), query, bookID, page, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responseList []response

	for rows.Next() {
		var review dbtypes.Review
		var user dbtypes.UserFromReview

		err = rows.Scan(
			&review.Id,
			&review.Title,
			&review.Content,
			&user.Id,
			&user.Name,
			&user.Profile_picture_url,
			&review.Total_likes,
			&review.Rating,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		// Append to responseList
		responseList = append(responseList, response{
			Review: review,
			User:   user,
		})
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		return nil, err
	}

	return responseList, nil
}

func SendReviewsById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	reviews, err := fetchBookReviews(conn, bookID, page, pageSize)
	if err != nil {
		log.Println("Error fetching reviews:", err)
		http.Error(w, "Error fetching reviews", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}