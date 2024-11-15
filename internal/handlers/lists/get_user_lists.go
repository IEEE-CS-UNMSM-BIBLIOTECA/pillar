package lists

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func GetUserLists(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	book_id := ps.ByName("id")
	username := r.Context().Value("username").(string)

	user_id := auth.GetIdFromUsername(username)
	if user_id == 0 {
		http.Error(w, "That username does not exist", http.StatusBadRequest)
		return
	}

	bookID, err := strconv.Atoi(book_id)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT * FROM get_user_lists($1, $2)`
	rows, err := conn.Query(context.Background(), query, user_id, bookID)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching user lists", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	lists := []dbtypes.ListAddDocument{}

	for rows.Next() {
		var list dbtypes.ListAddDocument
		err = rows.Scan(
			&list.ID,
			&list.Title,
			&list.HasDocument,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		lists = append(lists, list)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lists); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
