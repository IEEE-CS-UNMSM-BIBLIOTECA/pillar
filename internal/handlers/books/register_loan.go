package books

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"time"

	"github.com/julienschmidt/httprouter"
)

func RegisterLoan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var loan dbtypes.RegisterLoan
	if err := json.NewDecoder(r.Body).Decode(&loan); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	maxBooks := 5

	maxRetunDate, err := time.Parse("2006-01-02", loan.MaxRetunDate) // Expecting format YYYY-MM-DD
	if err != nil {
		log.Println("Invalid return date format:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `
		SELECT register_loan($1, $2, $3, $4);
	`, loan.BookID, loan.UserID, maxBooks, maxRetunDate)

	if err != nil {
		log.Println("Error calling register_loan:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Loan successfully registered"))
}
