package lists

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func GetAllLists(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	page := 1
	pageSize := 100

	queryParams := r.URL.Query()
	if p := queryParams.Get("page"); p != "" {
		parsedPage, err := strconv.Atoi(p)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if l := queryParams.Get("limit"); l != "" {
		parsedLimit, err := strconv.Atoi(l)
		if err == nil && parsedLimit > 0 {
			pageSize = parsedLimit
		}
	}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT * FROM get_all_lists($1, $2, $3)`
	rows, err := conn.Query(context.Background(), query, page, pageSize)
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
