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

func AddDocToList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list_id := ps.ByName("list_id")
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

	query := `SELECT add_book_to_list($1, $2)`
	rows, err := conn.Query(context.Background(), query, list_id, req.DocumentID)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error adding the book to the list", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Book added to the list successfully"}`))
}
