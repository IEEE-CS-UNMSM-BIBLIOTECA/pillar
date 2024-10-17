package handlers

import (
	"context"
	"log"
	"net/http"

	dbtypes "pillar/internal/db/types"
	dbutils "pillar/internal/db/utils"

    jsonexp "github.com/go-json-experiment/json"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

func getDocumentsBy(conn *pgxpool.Conn, field string, values []interface{}) []dbtypes.Document {
    query_str, query_params, query_err := goqu.
        Select("*").
        From("Document").
        Where(goqu.Ex{
            field: values,
        }).ToSQL()

    if query_err != nil {
        log.Println(query_err)
        return nil
    }

    rows, rows_err := conn.Query(context.Background(), query_str, query_params...)
    if rows_err != nil {
        log.Println(rows_err)
        return nil
    }


    for rows.Next() {
        scan_doc, scan_err := pgx.RowToStructByName[dbtypes.Document](rows)
        if scan_err != nil {
            log.Println(scan_err)
            rows.Close()
            return nil
        }

        log.Println(scan_doc)
    }


    log.Println("query:", query_str, "with params:", query_params)
    return nil
}

func HndGetDocumentsBy(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    type Req struct {
        Values []interface{} `json:"values"`
    }

    w.Header().Add("Access-Control-Allow-Origin", "*")
    w.Header().Add("Access-Control-Allow-Methods", "*")
    w.Header().Add("Access-Control-Allow-Headers", "*")

    new_req := Req {}
    field := params.ByName("field")
    if len(field) == 0 {
        log.Println("invalid fieldname")
        w.WriteHeader(http.StatusNotFound)
        return
    }

    unmarshal_err := jsonexp.UnmarshalRead(r.Body, &new_req, jsonexp.DefaultOptionsV2())
    if unmarshal_err != nil || new_req.Values == nil {
        log.Println(unmarshal_err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    conn, conn_err := dbutils.DbPool.Acquire(context.Background())
    if conn_err != nil {
        log.Println(conn_err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    defer conn.Release()

    getDocumentsBy(conn, field, new_req.Values)
}

