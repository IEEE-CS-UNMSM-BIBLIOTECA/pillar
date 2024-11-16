package dashboard

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"

	"github.com/julienschmidt/httprouter"
)

func AddPublisher(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var reviewReq dbtypes.Publisher
	if err := json.NewDecoder(r.Body).Decode(&reviewReq); err != nil {
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

	query := `INSERT INTO "Publisher" (name) VALUES ($1)`
	_, err = conn.Exec(context.Background(), query, reviewReq.Name)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error inserting the new publisher", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Publisher added successfully"}`))
}
