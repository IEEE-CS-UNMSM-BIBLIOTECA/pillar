package lists

import (
	"context"
	"log"
	"net/http"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"

	"github.com/julienschmidt/httprouter"
)

func CreateList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Context().Value("username").(string)

	user_id := auth.GetIdFromUsername(username)
	if user_id == 0 {
		http.Error(w, "That username does not exist", http.StatusBadRequest)
		return
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT create_empty_list($1)`
	rows, err := conn.Query(context.Background(), query, user_id)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error adding the book to the list", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Empty list created"}`))
}
