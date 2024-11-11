package handlers

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

func callToQuery(conn *pgxpool.Conn, ID int) (dbtypes.Document, dbtypes.Authors, dbtypes.Language, dbtypes.Publisher, dbtypes.Format, error) {
	var document dbtypes.Document
	var authors dbtypes.Authors
	var language dbtypes.Language
	var publisher dbtypes.Publisher
	var format dbtypes.Format
	var authorsJSON []byte

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
		&format.Id,
		&format.Name,
		&publisher.Id,
		&publisher.Name,
		&authorsJSON,
	)

	if err != nil {
		return dbtypes.Document{}, dbtypes.Authors{}, dbtypes.Language{}, dbtypes.Publisher{}, dbtypes.Format{}, err
	}

	err = json.Unmarshal(authorsJSON, &authors.Authors)
	if err != nil {
		log.Println("Error unmarshalling authors:", err)
		return dbtypes.Document{}, dbtypes.Authors{}, dbtypes.Language{}, dbtypes.Publisher{}, dbtypes.Format{}, err
	}

	return document, authors, language, publisher, format, nil
}

func sendBookById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	document, author, language, publisher, format, err := callToQuery(conn, bookID)
	if err != nil {
		log.Println("Error executing query:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Construct a combined response
	response := map[string]interface{}{
		"document":  document,
		"author":    author,
		"language":  language,
		"publisher": publisher,
		"format":    format,
	}

	// Set the response header to application/json and encode the result into JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
