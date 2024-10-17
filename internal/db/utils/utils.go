package utils

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/profclems/go-dotenv"
)

func newDbPool() *pgxpool.Pool {
    config_load_err := dotenv.Load()
    if config_load_err != nil {
        log.Panicln("couldn't load .env file, failed with:\n\t")
        log.Panicln(config_load_err)
    }

    conn_string := dotenv.Get("conn_string").(string)
    dbpool, pool_err := pgxpool.New(context.Background(), conn_string)
    if pool_err != nil {
        log.Panicln("pool generation error:\n\t", pool_err)
        return nil
    }

    return dbpool
}

var DbPool = newDbPool()


