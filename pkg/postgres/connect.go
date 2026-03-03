package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(dbDriver, dbUri string) *sqlx.DB {
	conn, err := sqlx.Open(dbDriver, dbUri)
	if err != nil {
		log.Fatalf("sqlx.Connect: %v", err)
	}

	if err := conn.PingContext(context.Background()); err != nil {
		log.Fatalf("sqlx.PingContext: %v", err)
	}

	conn.SetMaxIdleConns(15)
	conn.SetMaxOpenConns(25)
	conn.SetConnMaxLifetime(time.Minute * 5)

	return conn
}
