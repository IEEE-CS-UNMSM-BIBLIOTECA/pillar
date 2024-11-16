package images

import (
	"context"
	"log"
	"net/http"
	dbutils "pillar/internal/db/utils"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
)

func ImageLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("document_id")
	documentID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var coverURL []byte
	query := `SELECT cover_url FROM "Document" WHERE id = $1`
	err = conn.QueryRow(context.Background(), query, documentID).Scan(&coverURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
		log.Println("Error querying cover_url:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(coverURL) == 0 {
		http.Error(w, "No image available for this document", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(coverURL)
	if err != nil {
		log.Println("Error writing image to response:", err)
		http.Error(w, "Error sending image data", http.StatusInternalServerError)
	}
}
