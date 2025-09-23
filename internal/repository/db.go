package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewBunDB(ctx context.Context, connectionString string) (db *bun.DB, closeFunc func() error) {
	connector := pgdriver.NewConnector(pgdriver.WithDSN(connectionString))
	postgresql := sql.OpenDB(connector)
	db = bun.NewDB(postgresql, pgdialect.New(), bun.WithDiscardUnknownColumns())

	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("FATAL: %s (err=%s)\n", err, "could not ping postgres")
		os.Exit(1)
	}

	return db, func() error {
		return db.Close()
	}
}
