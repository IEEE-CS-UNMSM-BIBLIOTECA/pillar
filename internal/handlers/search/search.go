package search

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"

	"github.com/julienschmidt/httprouter"
)

func Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lookup := ps.ByName("lookup")
	if lookup == "" {
		http.Error(w, "Search term is required", http.StatusBadRequest)
		return
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
			d.id,
			d.title,
			COALESCE(
				(SELECT jsonb_agg(jsonb_build_object('id', a.id, 'name', a.name))
				FROM "Author_Document" ad
				JOIN "Author" a ON ad.author_id = a.id
				WHERE ad.document_id = d.id),
				'[]'::jsonb
			) AS authors
		FROM "Document" d
		LEFT JOIN "Tag_Document" td ON d.id = td.document_id
		LEFT JOIN "Tag" t ON td.tag_id = t.id
		WHERE d.title ILIKE '%' || $1 || '%'
		OR t.name ILIKE '%' || $1 || '%'
		GROUP BY d.id, d.title, d.cover_url;
	`

	rows, err := conn.Query(context.Background(), query, lookup)
	if err != nil {
		log.Println("Query execution failed:", err)
		http.Error(w, "Error executing search query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var id int
		var title *string
		var authorsJSON []byte

		err = rows.Scan(&id, &title, &authorsJSON)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing result", http.StatusInternalServerError)
			return
		}

		var authors dbtypes.Authors

		log.Printf("Authors JSON: %s\n", authorsJSON)
		if err := json.Unmarshal(authorsJSON, &authors.Authors); err != nil {
			log.Println("Failed to parse authors:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		urlImage := fmt.Sprintf("http://143.198.142.139:8080/cover/%d", id)

		result := map[string]interface{}{
			"id":        id,
			"title":     title,
			"authors":   authors.Authors,
			"cover_url": urlImage,
		}

		results = append(results, result)
	}

	if len(results) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
