package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"time"

	"github.com/julienschmidt/httprouter"
)

func GetLanguages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var languages = []dbtypes.Language{}

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
	var publishers = []dbtypes.Publisher{}

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
	var countries = []dbtypes.Country{}

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
	var formats = []dbtypes.Format{}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT * FROM "DocumentFormat"`
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
	var authors = []dbtypes.Author{}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT id, name FROM "Author"`
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

func GetGenders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var genders = []dbtypes.Gender{}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT id, name FROM "Gender"`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching genders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var gender dbtypes.Gender
		err = rows.Scan(
			&gender.Id,
			&gender.Name,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		genders = append(genders, gender)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(genders); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetDocuments(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var documents = []dbtypes.Document{}

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
			d.isbn, 
			d.description, 
			d.publication_year, 
			d.acquisition_date, 
			d.edition, 
			d.external_lend_allowed, 
			d.total_pages, 
			d.base_price, 
			d.total_copies, 
			d.available_copies, 
			COALESCE(d.mean_rating, 0) AS mean_rating,
			-- Authors
			(SELECT jsonb_agg(jsonb_build_object(
				'id', a.id, 
				'name', a.name
			)) FROM "Author" a 
			JOIN "Author_Document" da ON da.author_id = a.id 
			WHERE da.document_id = d.id) AS authors,
			-- Language as a JSON object
			jsonb_build_object(
				'id', l.id, 
				'name', l.name
			) AS language,
			-- Publisher as a JSON object
			jsonb_build_object(
				'id', p.id, 
				'name', p.name
			) AS publisher,
			-- Tags as a JSONB array
			(SELECT jsonb_agg(jsonb_build_object(
				'id', t.id, 
				'name', t.name
			)) FROM "Tag" t 
			JOIN "Tag_Document" dt ON dt.tag_id = t.id 
			WHERE dt.document_id = d.id) AS tags,
			-- Document format as a JSONB object
			jsonb_build_object(
				'id', f.id, 
				'name', f.name
			) AS formats
		FROM 
			"Document" d
		JOIN 
			"Language" l ON l.id = d.language_id
		JOIN 
			"Publisher" p ON p.id = d.publisher_id
		LEFT JOIN 
			"DocumentFormat" f ON f.id = d.format_id;
	`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching documents", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var document dbtypes.Document
		var acquisition_date time.Time

		var authorJSON []byte
		var languageJSON []byte
		var formatJSON []byte
		var publisherJSON []byte
		var tagJSON []byte

		err = rows.Scan(
			&document.Id,
			&document.Title,
			&document.Isbn,
			&document.Description,
			&document.Publication_year,
			&acquisition_date,
			&document.Edition,
			&document.External_lend_allowed,
			&document.Total_pages,
			&document.Base_price,
			&document.Total_copies,
			&document.Available_copies,
			&document.Mean_rating,
			&authorJSON,
			&languageJSON,
			&publisherJSON,
			&tagJSON,
			&formatJSON,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}

		var authors []dbtypes.Author
		if err := json.Unmarshal(authorJSON, &authors); err != nil {
			log.Println("Error unmarshalling authors:", err)
			http.Error(w, "Error processing authors", http.StatusInternalServerError)
			return
		}
		document.Authors = authors

		var language dbtypes.Language
		if err := json.Unmarshal(languageJSON, &language); err != nil {
			log.Println("Error unmarshalling language:", err)
			http.Error(w, "Error processing language", http.StatusInternalServerError)
			return
		}
		document.Language = language

		var formats dbtypes.Format
		if err := json.Unmarshal(formatJSON, &formats); err != nil {
			log.Println("Error unmarshalling formats:", err)
			http.Error(w, "Error processing formats", http.StatusInternalServerError)
			return
		}
		document.Format = formats

		var publisher dbtypes.Publisher
		if err := json.Unmarshal(publisherJSON, &publisher); err != nil {
			log.Println("Error unmarshalling publisher:", err)
			http.Error(w, "Error processing publisher", http.StatusInternalServerError)
			return
		}
		document.Publisher = publisher

		var tag []dbtypes.Tag
		if err := json.Unmarshal(tagJSON, &tag); err != nil {
			log.Println("Error unmarshalling tags:", err)
			http.Error(w, "Error processing tags", http.StatusInternalServerError)
			return
		}
		document.Tags = tag

		document.Acquisition_date = acquisition_date.Format("2006-01-02")

		urlImage := fmt.Sprintf("http://143.198.142.139:8080/cover/%d", document.Id)
		document.Cover_url = &urlImage

		documents = append(documents, document)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	// Set response header to application/json and encode the map
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(documents); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetOrders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var orders = []dbtypes.OrderView{}

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `
		SELECT 
			o.id AS order_id, 
			o.order_date, 
			o.max_return_date, 
			o.actual_return_date,
			u.id AS user_id,
			u.username AS name,
			d.id AS document_id,
			d.title AS document_title
		FROM 
			"Order" o
		JOIN 
			"User" u ON o.user_id = u.id
		JOIN 
			"Document" d ON o.document_id = d.id
			`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order dbtypes.OrderView

		var orderDate time.Time
		var maxReturnDate time.Time
		var actualReturnDate *time.Time

		err = rows.Scan(
			&order.Id,
			&orderDate,
			&maxReturnDate,
			&actualReturnDate,
			&order.User.Id,
			&order.User.Name,
			&order.Document.Id,
			&order.Document.Title,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}

		order.Order_date = orderDate.Format("2006-01-02")
		order.Max_return_date = maxReturnDate.Format("2006-01-02")
		if actualReturnDate != nil {
			order.Actual_return_date = actualReturnDate.Format("2006-01-02")
		} else {
			order.Actual_return_date = ""
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	// Set response header to application/json and encode the map
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetAuthorsByDoc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	document_id := ps.ByName("document_id")

	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := `SELECT a.id, a.name
	FROM "Author" a 
	JOIN "Author_Document" ad ON a.id = ad.author_id
	JOIN "Document" d ON ad.document_id = d.id
	WHERE d.id =$1`
	rows, err := conn.Query(context.Background(), query, document_id)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Error fetching orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var authors = []dbtypes.Author{}
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

	// Set response header to application/json and encode the map
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authors); err != nil {
		log.Println("Error encoding response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
