package books

import (
	"context"
	"log"
	"net/http"
	dbutils "pillar/internal/db/utils"

	"github.com/julienschmidt/httprouter"
)

func AddReviews(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// document_id := ps.ByName("document_id")
	// user_id := ps.ByName("user_id")
	// bookID, err := strconv.Atoi(document_id)
	// if err != nil {
	// 	http.Error(w, "Invalid book ID", http.StatusBadRequest)
	// 	return
	// }

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

}
