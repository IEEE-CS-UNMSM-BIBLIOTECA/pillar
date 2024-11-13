package lists

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"

	"github.com/julienschmidt/httprouter"
)

func DeleteDocFromList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req dbtypes.AddDocList
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT remove_book_from_list($1, $2)`
	rows, err := conn.Query(context.Background(), query, req.ListID, req.DocumentID)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error adding the book to the list", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Book deleted from the list successfully"}`))
}
