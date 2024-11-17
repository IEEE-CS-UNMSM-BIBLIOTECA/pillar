package books

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

func GetBooksFromList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list_id := ps.ByName("list_id")

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `
	SELECT * FROM get_documents_by_list($1);
	`
	rows, err := conn.Query(context.Background(), query, list_id)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching user lists", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var documents = []dbtypes.PopularBook{}
	var authorsJSON []byte

	for rows.Next() {
		var document dbtypes.PopularBook
		err = rows.Scan(
			&document.BookID,
			&document.Title,
			&authorsJSON,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(authorsJSON, &document.Authors)
		if err != nil {
			log.Println("Error unmarshaling authors:", err)
			http.Error(w, "Error processing authors", http.StatusInternalServerError)
			return
		}

		urlImage := fmt.Sprintf("http://143.198.142.139:8080/cover/%d", document.BookID)
		document.CoverURL = &urlImage

		documents = append(documents, document)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(documents); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
