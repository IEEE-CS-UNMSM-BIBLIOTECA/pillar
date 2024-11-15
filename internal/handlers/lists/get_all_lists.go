package lists

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"

	"github.com/julienschmidt/httprouter"
)

func GetAllLists(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req struct {
		PageSize int    `json:"page_size"`
		Page     int    `json:"page"`
		Tags     *[]int `json:"tags"`
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT * FROM get_all_lists($1, $2, $3)`
	rows, err := conn.Query(context.Background(), query, req.PageSize, req.Page, req.Tags)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching user lists", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lists []dbtypes.ListAddDocument

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
