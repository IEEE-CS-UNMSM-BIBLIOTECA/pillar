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

func GetAllLists(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username, hasToken := r.Context().Value("username").(string)

	var user_id int
	if hasToken {
		user_id = auth.GetIdFromUsername(username)
		if user_id == 0 {
			http.Error(w, "That username does not exist", http.StatusBadRequest)
			return
		}
	}

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

	query := `SELECT * FROM get_all_lists($1, $2, $3, $4)`
	rows, err := conn.Query(context.Background(), query, user_id, pageSize, page, nil)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching user lists", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lists = []dbtypes.List{}

	for rows.Next() {
		var list dbtypes.List
		err = rows.Scan(
			&list.Id,
			&list.Title,
			&list.Total_likes,
			&list.Total_books,
			&list.Preview_images,
			&list.Private,
			&list.Liked,
			&list.Own,
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
