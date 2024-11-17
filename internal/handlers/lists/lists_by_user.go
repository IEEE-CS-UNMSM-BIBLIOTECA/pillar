package lists

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"

	"github.com/julienschmidt/httprouter"
)

func GetListByUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// user_lookup := ps.ByName("user_id")

	username, hasToken := r.Context().Value("username").(string)
	var user_id int
	if hasToken {
		user_id = auth.GetIdFromUsername(username)
		if user_id == 0 {
			http.Error(w, "That username does not exist", http.StatusBadRequest)
			return
		}
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	var list dbtypes.List
	query := `
	SELECT * FROM "List"
	JOIN "ListLike" lk ON 
	WHERE user_id = $1 AND private = false
	`
	err = conn.QueryRow(context.Background(), query, user_id).Scan(
		&list.Id,
		&list.Title,
		&list.Total_likes,
		&list.Total_books,
		&list.Liked,
		&list.Own,
		&list.User.Id,
		&list.User.Name,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			http.Error(w, "List not found", http.StatusNotFound)
		} else {
			log.Println("Error executing query:", err)
			http.Error(w, "Error fetching list", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
