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

func AddTag(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var tagReq dbtypes.Tag
	if err := json.NewDecoder(r.Body).Decode(&tagReq); err != nil {
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

	query := `INSERT INTO "Tag" (name) VALUES ($1)`
	_, err = conn.Exec(context.Background(), query, tagReq.Name)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error inserting the new tag", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Tag added successfully"}`))
}
