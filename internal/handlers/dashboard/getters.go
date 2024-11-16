package dashboard

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"

	"github.com/julienschmidt/httprouter"
)

func GetLanguages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var languages []dbtypes.Language

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT * FROM "Language"`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching languages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var list dbtypes.Language
		err = rows.Scan(
			&list.Id,
			&list.Name,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		languages = append(languages, list)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(languages); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetPublishers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var publishers []dbtypes.Publisher

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT * FROM "Publisher"`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching languages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var publisher dbtypes.Publisher
		err = rows.Scan(
			&publisher.Id,
			&publisher.Name,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		publishers = append(publishers, publisher)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(publishers); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}