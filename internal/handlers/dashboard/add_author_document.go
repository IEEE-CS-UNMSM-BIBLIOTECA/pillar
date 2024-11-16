package dashboard

import (
	"context"
	"log"
	"net/http"
	dbutils "pillar/internal/db/utils"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func AddAuthorDocument(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	author_id := ps.ByName("author_id")
	authorID, err := strconv.Atoi(author_id)
	if err != nil {
		http.Error(w, "Invalid document id", http.StatusBadRequest)
		return
	}

	document_id := ps.ByName("document_id")
	documentID, err := strconv.Atoi(document_id)
	if err != nil {
		http.Error(w, "Invalid document id", http.StatusBadRequest)
		return
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `INSERT INTO "Author_Document" (author_id, document_id) VALUES ($1, $2)`
	_, err = conn.Exec(context.Background(), query, authorID, documentID)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error inserting the document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Author_document added successfully"}`))
}
