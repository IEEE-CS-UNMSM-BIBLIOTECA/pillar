package search

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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
		SELECT DISTINCT d.id,
				d."title",
				a."name" AS autor,
				a.id AS autor_id,
				d.cover_url
		FROM "Document" d
		LEFT JOIN "Author_Document" ad ON d.id = ad.document_id
		LEFT JOIN "Author" a ON ad.author_id = a.id
		LEFT JOIN "Tag_Document" td ON d.id = td.document_id
		LEFT JOIN "Tag" t ON td.tag_id = t.id
		WHERE d."title" ILIKE '%' || $1 || '%'
		OR t."name" ILIKE '%' || $1 || '%';
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
		var title, authorName, coverUrl *string
		var authorID *int

		err = rows.Scan(&id, &title, &authorName, &authorID, &coverUrl)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing result", http.StatusInternalServerError)
			return
		}

		result := map[string]interface{}{
			"id":        id,
			"title":     title,
			"author":    authorName,
			"author_id": authorID,
			"cover_url": coverUrl,
		}

		results = append(results, result)
	}

	if len(results) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "No documents found"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
