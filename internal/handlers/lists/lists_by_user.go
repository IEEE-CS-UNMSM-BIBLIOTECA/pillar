package lists

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"

	"github.com/julienschmidt/httprouter"
)

func GetListByUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_lookup := ps.ByName("user_id")

	username, hasToken := r.Context().Value("username").(string)
	var requesterID int
	if hasToken {
		requesterID = auth.GetIdFromUsername(username)
		if requesterID == 0 {
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

	query := `
	SELECT
		l.id,
		l.title,
		l.total_likes,
		l.total_books,
		CASE WHEN $1 = l.user_id THEN l.private ELSE false END AS is_private,
		EXISTS(SELECT 1 FROM "ListLike" lk WHERE lk.list_id = l.id AND lk.user_id = $1) AS liked,
		l.user_id = $1 AS own,
		ARRAY(
			SELECT d.id
			FROM "List_Document" ld
			JOIN "Document" d ON ld.document_id = d.id
			WHERE ld.list_id = l.id
			LIMIT 5
		) AS document_ids
	FROM "List" l
	JOIN "User" u ON l.user_id = u.id
	WHERE l.user_id = $2
	`
	rows, err := conn.Query(context.Background(), query, requesterID, user_lookup)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching lists", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lists = []dbtypes.ListUser{}
	var documentIDs []int

	for rows.Next() {
		var list dbtypes.ListUser
		err = rows.Scan(
			&list.Id,
			&list.Title,
			&list.Total_likes,
			&list.Total_books,
			&list.Private,
			&list.Liked,
			&list.Own,
			&documentIDs,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing lists", http.StatusInternalServerError)
			return
		}
		for _, docID := range documentIDs {
			list.Preview_images = append(list.Preview_images, fmt.Sprintf("http://143.198.142.139:8080/cover/%d", docID))
		}

		lists = append(lists, list)
	}

	if len(lists) == 0 {
		http.Error(w, "No lists found", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lists); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}
