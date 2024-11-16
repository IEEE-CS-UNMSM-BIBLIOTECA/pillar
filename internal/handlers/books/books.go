package books

import (
	"context"
	"fmt"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"
	"strconv"

	"github.com/go-json-experiment/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

func getPopularTags(conn *pgxpool.Conn) ([]int, error) {
	query := `
        SELECT id
        FROM "Tag"
        ORDER BY mean_rating DESC
        LIMIT 5;
    `

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular tags: %w", err)
	}
	defer rows.Close()

	var popularTags []int
	for rows.Next() {
		var tagID int
		if err := rows.Scan(&tagID); err != nil {
			return nil, fmt.Errorf("failed to scan tag ID: %w", err)
		}
		popularTags = append(popularTags, tagID)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating rows: %w", rows.Err())
	}

	return popularTags, nil
}

func getUserTopTags(conn *pgxpool.Conn, userID, topLimit int) ([]int, error) {
	rows, err := conn.Query(context.Background(), `SELECT tag_id FROM get_user_top_tags($1, $2)`, userID, topLimit)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tagIDs []int

	for rows.Next() {
		var tagID int
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
	rows, err := conn.Query(context.Background(), `
        SELECT book_id, title, authors
        FROM get_popular_books($1, $2, $3)`, page, pageSize, tagIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []dbtypes.PopularBook
	for rows.Next() {
		var book dbtypes.PopularBook
		var authorsJSON []byte

		if err := rows.Scan(&book.BookID, &book.Title, &authorsJSON); err != nil {
			return nil, err
		}

		var authors []dbtypes.Author
		if err := json.Unmarshal(authorsJSON, &authors); err != nil {
			log.Println("Error unmarshalling authors:", err)
			return nil, err
		}
		urlImage := fmt.Sprintf("http://143.198.142.139:8080/cover/%d", book.BookID)

		book.Authors = authors
		book.CoverURL = &urlImage

		books = append(books, book)
	}

	return books, nil
}

func SendPopularBooks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.Context().Value("username").(string)

	user_id := auth.GetIdFromUsername(username)
	if user_id == 0 {
		http.Error(w, "That username does not exist", http.StatusBadRequest)
		return
	}

	page := 1
	pageSize := 100

	// Parse query parameters
	queryParams := r.URL.Query()
	if p := queryParams.Get("page"); p != "" {
		parsedPage, err := strconv.Atoi(p)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if l := queryParams.Get("limit"); l != "" {
		parsedLimit, err := strconv.Atoi(l)
		if err == nil && parsedLimit > 0 {
			pageSize = parsedLimit
		}
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
	topTags, err := getUserTopTags(conn, user_id, 5)
	if err != nil {
		log.Println("Failed to get top tags:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var books []dbtypes.PopularBook

	if len(topTags) == 0 {
		log.Println("Calling popular tags instead")
		popularTags, err := getPopularTags(conn)
		if err != nil {
			log.Println("Failed to get popular tags:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Retrieved popular tags: %v", popularTags)
		books, err = getPopularBooks(conn, page, pageSize, popularTags)
		if err != nil {
			log.Println("Failed to get popular books:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		books, err = getPopularBooks(conn, page, pageSize, topTags)
		if err != nil {
			log.Println("Failed to get popular books:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
