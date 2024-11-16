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

func AddAuthorDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var reviewReq dbtypes.AuthorDashboard
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

	query := `INSERT INTO "Author_Document" (author_id, birth_date) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = conn.Exec(context.Background(), query, reviewReq.Name, reviewReq.Bio, reviewReq.GenderID,
		reviewReq.CountryID, reviewReq.ImageUrl)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error inserting the document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Author added successfully"}`))
}
