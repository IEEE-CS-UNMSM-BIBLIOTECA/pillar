package books

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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

func callToQuery(conn *pgxpool.Conn, ID int) (dbtypes.Document, dbtypes.Authors, dbtypes.Language, dbtypes.Publisher, dbtypes.Tags, error) {
	var document dbtypes.Document
	var authors dbtypes.Authors
	var tags dbtypes.Tags
	var language dbtypes.Language
	var publisher dbtypes.Publisher
	var authorsJSON []byte
	var tagsJson []byte

	var acquisition_date pgtype.Date

	query := `SELECT * FROM get_book_details($1)`
	row := conn.QueryRow(context.Background(), query, ID)

	err := row.Scan(
		&document.Id,
		&document.Title,
		&document.Isbn,
		&document.Description,
		&document.Cover_url,
		&document.Publication_year,
		&acquisition_date,
		&document.Edition,
		&document.Total_pages,
		&document.External_lend_allowed,
		&document.Base_price,
		&document.Total_copies,
		&document.Available_copies,
		&document.Mean_rating,
		&language.Id,
		&language.Name,
		&publisher.Id,
		&publisher.Name,
		&authorsJSON,
		&tagsJson,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return dbtypes.Document{}, dbtypes.Authors{}, dbtypes.Language{}, dbtypes.Publisher{}, dbtypes.Tags{}, err
		}
		log.Println("Error scanning row:", err)
		return dbtypes.Document{}, dbtypes.Authors{}, dbtypes.Language{}, dbtypes.Publisher{}, dbtypes.Tags{}, err
	}

	document.Acquisition_date = acquisition_date.Time.Format("2006-01-02")

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

	urlImage := fmt.Sprintf("http://143.198.142.139:8080/cover/%d", bookID)

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	document, author, language, publisher, tag, err := callToQuery(conn, bookID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Book with ID %d not found\n", bookID)
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		log.Println("Error executing query:", err)
		http.Error(w, "Error retrieving book details", http.StatusInternalServerError)
		return
	}

	documentMap := map[string]interface{}{
		"mean_rating":           document.Mean_rating,
		"id":                    bookID,
		"publication_year":      document.Publication_year,
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
		"cover_url":             urlImage,
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

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
