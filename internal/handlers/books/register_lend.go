package books

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"
	"time"

	"github.com/julienschmidt/httprouter"
)

func RegisterOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Context().Value("username").(string)

	user_id := auth.GetIdFromUsername(username)
	if user_id == 0 {
		http.Error(w, "That username does not exist", http.StatusBadRequest)
		return
	}

	var lend dbtypes.RegisterLend
	if err := json.NewDecoder(r.Body).Decode(&lend); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	maxBooks := 5

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	maxReturnDate := time.Now().AddDate(0, 0, 7)

	_, err = conn.Exec(context.Background(), `
		SELECT register_loan($1, $2, $3, $4);
	`, lend.BookID, user_id, maxBooks, maxReturnDate)

	if err != nil {
		log.Println("Error calling register_loan:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Loan successfully registered"))
}
