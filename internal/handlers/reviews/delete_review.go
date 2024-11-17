package reviews

import (
	"context"
	"log"
	"net/http"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func DeleteReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	review := ps.ByName("review_id")

	username, hasToken := r.Context().Value("username").(string)
	var requesterID int
	if hasToken {
		requesterID = auth.GetIdFromUsername(username)
		if requesterID == 0 {
			http.Error(w, "That username does not exist", http.StatusBadRequest)
			return
		}
	}

	reviewID, err := strconv.Atoi(review)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var reviewOwnerID int
	query := `SELECT user_id FROM "Review" WHERE id = $1`
	err = conn.QueryRow(context.Background(), query, reviewID).Scan(&reviewOwnerID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			http.Error(w, "Review not found", http.StatusNotFound)
		} else {
			log.Println("Error checking review owner:", err)
			http.Error(w, "Error fetching review", http.StatusInternalServerError)
		}
		return
	}

	if reviewOwnerID != requesterID {
		http.Error(w, "You do not have permission to delete this review", http.StatusForbidden)
		return
	}

	deleteQuery := `DELETE FROM "Review" WHERE id = $1`
	_, err = conn.Exec(context.Background(), deleteQuery, reviewID)
	if err != nil {
		log.Println("Error deleting review:", err)
		http.Error(w, "Error deleting review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Review deleted successfully"))
}
