package lends

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	dbutypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"
	"pillar/internal/handlers/auth"
	"time"

	"github.com/julienschmidt/httprouter"
)

func StringPtr(s string) *string {
	return &s
}

func GetLendsByUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Context().Value("username").(string)

	user_id := auth.GetIdFromUsername(username)
	if user_id == 0 {
		http.Error(w, "That username does not exist", http.StatusBadRequest)
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
            o.id AS order_id,
            o.order_date,
            o.max_return_date,
            o.actual_return_date,
            d.id AS document_id,
            d.title AS document_title,
            COALESCE(
			(SELECT jsonb_agg(jsonb_build_object('id', a.id, 'name', a.name))
			FROM "Author_Document" ad
			JOIN "Author" a ON a.id = ad.author_id
			WHERE ad.document_id = d.id),
			'[]'::jsonb
		) AS authors
        FROM "Order" o
        JOIN "Document" d ON o.document_id = d.id
        LEFT JOIN "Author_Document" da ON d.id = da.document_id
        LEFT JOIN "Author" a ON da.author_id = a.id
        WHERE o.user_id = $1 AND o.actual_return_date IS NULL
        GROUP BY o.id, d.id;
    `

	rows, err := conn.Query(context.Background(), query, user_id)
	if err != nil {
		log.Println("Failed to execute query:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	orders := []dbutypes.OrderRequest{}
	for rows.Next() {
		var order dbutypes.OrderRequest
		var authors []byte

		var max_return_date time.Time
		var order_date time.Time
		var actual_return_date *time.Time

		if err := rows.Scan(
			&order.Id,
			&order_date,
			&max_return_date,
			&actual_return_date,
			&order.Document.BookID,
			&order.Document.Title,
			&authors,
		); err != nil {
			log.Println("Failed to scan row:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if actual_return_date != nil {
			order.Actual_return_date = StringPtr(actual_return_date.Format("2006-01-02"))
		} else {
			order.Actual_return_date = nil
		}

		order.Max_return_date = max_return_date.Format("2006-01-02")
		order.Order_date = order_date.Format("2006-01-02")

		urlImage := fmt.Sprintf("http://143.198.142.139:8080/cover/%d", order.Document.BookID)

		order.Document.CoverURL = &urlImage

		if err := json.Unmarshal([]byte(authors), &order.Document.Authors); err != nil {
			log.Println("Failed to parse authors:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
