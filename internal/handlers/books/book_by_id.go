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

func callToQuery(conn *pgxpool.Conn, ID int) (dbtypes.Document, dbtypes.Authors, dbtypes.Language, dbtypes.Publisher, dbtypes.Tags, error) {
	var document dbtypes.Document
	var authors dbtypes.Authors
	var language dbtypes.Language
	var publisher dbtypes.Publisher
	var tags dbtypes.Tags
	var authorsJSON []byte
	var tagsJson []byte

	query := `SELECT * FROM get_book_details($1)`
	row := conn.QueryRow(context.Background(), query, ID)

	err := row.Scan(
		&document.Id,
		&document.Title,
		&document.Isbn,
		&document.Description,
		&document.Cover_url,
		&document.Publication_date,
		&document.Acquisition_date,
		&document.Edition,
		&document.Total_pages,
		&document.External_lend_allowed,
		&document.Base_price,
		&document.Total_copies,
		&document.Available_copies,
		&document.Avg_rating,
		&language.Id,
		&language.Name,
		&publisher.Id,
		&publisher.Name,
		&authorsJSON,
		&tagsJson,
	)

	if err != nil {
		return dbtypes.Document{}, dbtypes.Authors{}, dbtypes.Language{}, dbtypes.Publisher{}, dbtypes.Tags{}, err
	}

	err = json.Unmarshal(authorsJSON, &authors.Authors)
	if err != nil {
		log.Println("Error unmarshalling authors:", err)
		return dbtypes.Document{}, dbtypes.Authors{}, dbtypes.Language{}, dbtypes.Publisher{}, dbtypes.Tags{}, err
	}

	err = json.Unmarshal(tagsJson, &tags.Tags)
	if err != nil {
		log.Println("Error unmarshalling tags:", err)
		return dbtypes.Document{}, dbtypes.Authors{}, dbtypes.Language{}, dbtypes.Publisher{}, dbtypes.Tags{}, err
	}

	return document, authors, language, publisher, tags, nil
}

func SendBookById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	document, author, language, publisher, tag, err := callToQuery(conn, bookID)
	if err != nil {
		log.Println("Error executing query:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	documentMap := map[string]interface{}{
		"avg_rating":            document.Avg_rating,
		"id":                    document.Id,
		"publication_date":      document.Publication_date,
		"acquisition_date":      document.Acquisition_date,
		"edition":               document.Edition,
		"total_pages":           document.Total_pages,
		"external_lend_allowed": document.External_lend_allowed,
		"base_price":            document.Base_price,
		"total_copies":          document.Total_copies,
		"available_copies":      document.Available_copies,
		"title":                 document.Title,
		"isbn":                  document.Isbn,
		"description":           document.Description,
		"cover_url":             document.Cover_url,
	}

	// Construct a combined response
	response := map[string]interface{}{
		"authors":   author.Authors,
		"language":  language,
		"publisher": publisher,
		"tags":      tag.Tags,
	}

	for key, value := range documentMap {
		response[key] = value
	}

	// Set the response header to application/json and encode the result into JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
