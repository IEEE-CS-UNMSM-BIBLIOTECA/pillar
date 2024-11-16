package dashboard

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"time"

	"github.com/julienschmidt/httprouter"
)

func AddDocToDB(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var reviewReq dbtypes.Document
	if err := json.NewDecoder(r.Body).Decode(&reviewReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	acquisition_date, err := time.Parse("2006-01-02", reviewReq.Acquisition_date) // Expecting format YYYY-MM-DD
	if err != nil {
		log.Println("Invalid birth date format:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `INSERT INTO "Document" (title, isbn, description, cover_url, acquisition_date, edition, total_pages, 
	external_lend_allowed, base_price, total_copies, available_copies, language_id, format_id, publisher_id, mean_rating, publication_year) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	_, err = conn.Exec(context.Background(), query, reviewReq.Title, reviewReq.Isbn, reviewReq.Description, reviewReq.Cover_url, acquisition_date,
		reviewReq.Edition, reviewReq.Total_pages, reviewReq.External_lend_allowed, reviewReq.Base_price, reviewReq.Total_copies, reviewReq.Available_copies,
		reviewReq.Language_id, reviewReq.Format_id, reviewReq.Publisher_id, reviewReq.Mean_rating, reviewReq.Publication_date)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error inserting the document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Document added successfully"}`))
}
