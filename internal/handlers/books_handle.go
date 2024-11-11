package handlers

import (
	"context"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"

	"github.com/doug-martin/goqu/v9"
	"github.com/go-json-experiment/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

func getUserTopTags(conn *pgxpool.Conn, userID, topLimit int) ([]int, error) {
	query_str, _, query_err := goqu.
		Select("*").
		From("Tag").ToSQL()

	if query_err != nil {
		log.Println(query_err)
		return nil, nil
	}

	rows, err := conn.Query(context.Background(), query_str)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tagIDs []int

	for rows.Next() {
		var tagID int
		log.Printf("ALGO:")
		if err := rows.Scan(&tagID); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		tagIDs = append(tagIDs, tagID)
	}

	log.Printf("Retrieved Tags: %v", tagIDs)

	return tagIDs, nil
}

func getPopularBooks(conn *pgxpool.Conn, page, pageSize int, tagIDs []int) ([]dbtypes.PopularBook, error) {
	if tagIDs == nil {
		return nil, nil
	}

	rows, err := conn.Query(context.Background(), `
        SELECT book_id, title, author_id, author_name, cover_url 
        FROM get_popular_books($1, $2, $3)`, page, pageSize, tagIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []dbtypes.PopularBook
	for rows.Next() {
		var book dbtypes.PopularBook
		if err := rows.Scan(&book.BookID, &book.Title, &book.AuthorID, &book.AuthorName, &book.CoverURL); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func sendPopularBooks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req struct {
		UserID   int `json:"user_id"`
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := json.UnmarshalRead(r.Body, &req, json.DefaultOptionsV2()); err != nil {
		log.Println("Failed to parse request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Acquire a database connection from the pool
	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	// Retrieve user's top tags
	// topTags, err := getUserTopTags(conn, req.UserID, 5)
	// if err != nil {
	// 	log.Println("Failed to get top tags:", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	query_str, _, query_err := goqu.
		Select("*").
		From("Tag").ToSQL()

	if query_err != nil {
		log.Println(query_err)
		return
	}

	rows, err := conn.Query(context.Background(), query_str)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	var tagIDs []int

	for rows.Next() {
		var tagID int
		log.Printf("ALGO:")
		if err := rows.Scan(&tagID); err != nil {
			log.Printf("Error scanning row: %v", err)
			return
		}

		tagIDs = append(tagIDs, tagID)
	}

	log.Printf("Retrieved Tags: %v", tagIDs)

	// Retrieve popular books using the top tags
	books, err := getPopularBooks(conn, req.Page, req.PageSize, tagIDs)
	if err != nil {
		log.Println("Failed to get popular books:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(books)
	if err != nil {
		log.Println("Failed to encode response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
