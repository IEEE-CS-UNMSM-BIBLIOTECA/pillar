package dashboard

import (
	"context"
	"log"
	"net/http"

	"pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"

	"github.com/doug-martin/goqu/v9"
	jsonexp "github.com/go-json-experiment/json"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
)

func GetAllTags(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conn, err := dbutils.DbPool.Acquire(context.Background())
	if err != nil {
		log.Println("Failed to acquire a database connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query_str, query_params, query_err := goqu.Select("*").From("Tag").ToSQL()
	if query_err != nil {
		log.Println(query_err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rows, row_query_err := conn.Query(context.Background(), query_str, query_params...)
	if row_query_err != nil {
		log.Println(row_query_err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tag_slice, tag_slice_err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Tag])
	if tag_slice_err != nil {
		log.Println(tag_slice_err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonexp.MarshalWrite(w, tag_slice, jsonexp.DefaultOptionsV2())
}
