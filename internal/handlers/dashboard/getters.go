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

func GetCountries(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var countries []dbtypes.Country

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT * FROM "Country"`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching languages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var countrie dbtypes.Country
		err = rows.Scan(
			&countrie.Id,
			&countrie.Name,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		countries = append(countries, countrie)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(countries); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetFormats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var formats []dbtypes.Format

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT * FROM "Format"`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching formats", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var format dbtypes.Format
		err = rows.Scan(
			&format.Id,
			&format.Name,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		formats = append(formats, format)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(formats); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetAuhors(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var authors []dbtypes.Author

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT * FROM "Author"`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching formats", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var author dbtypes.Author
		err = rows.Scan(
			&author.Id,
			&author.Name,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		authors = append(authors, author)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authors); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
