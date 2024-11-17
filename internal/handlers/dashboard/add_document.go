package dashboard

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"time"

	"github.com/julienschmidt/httprouter"
)

func AddDocToDB(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var documentReq dbtypes.Document
	if err := json.NewDecoder(r.Body).Decode(&documentReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	acquisitionDate, err := time.Parse("2006-01-02", documentReq.Acquisition_date)
	if err != nil {
		log.Println("Invalid birth date format:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var coverBytes []byte
	if documentReq.Cover_url != nil {
		coverBytes, err = base64.StdEncoding.DecodeString(*documentReq.Cover_url)
		if err != nil {
			log.Println("Invalid base64 in cover_url:", err)
			http.Error(w, "Invalid cover_url format", http.StatusBadRequest)
			return
		}
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var newDocumentID int
	query := `
		SELECT create_document(
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)`
	err = conn.QueryRow(context.Background(),
		query,
		documentReq.Title,
		documentReq.Isbn,
		documentReq.Description,
		acquisitionDate,
		documentReq.Edition,
		documentReq.Total_pages,
		documentReq.External_lend_allowed,
		documentReq.Base_price,
		documentReq.Total_copies,
		documentReq.Available_copies,
		documentReq.Language_id,
		documentReq.Format_id,
		documentReq.Publisher_id,
		documentReq.Publication_year,
		documentReq.Authors_id,
		documentReq.Tags_id,
		coverBytes,
	).Scan(&newDocumentID)
	if err != nil {
		log.Println("Error executing create_document function:", err)
		http.Error(w, "Error inserting the document", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Document added successfully",
		"book_id": newDocumentID,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}
}
