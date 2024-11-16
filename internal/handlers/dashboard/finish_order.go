package dashboard

import (
	"context"
	"log"
	"net/http"
	dbutils "pillar/internal/db/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
)

func FinishOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("order_id")

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	actualReturnDate := time.Now()

	query := `UPDATE "Order" SET actual_return_date = $1 WHERE id = $2`
	_, err = conn.Exec(context.Background(), query, actualReturnDate, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("No order found with that id")
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Println("Error executing query:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order updated correctly"))
}
