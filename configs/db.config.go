package configs

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB_POOL *pgxpool.Pool

func InitDBCon() error {
	DSN := os.Getenv("DSN")

	db, err := pgxpool.New(context.Background(), DSN)
	if err != nil {
		return err
	}

	if err := db.Ping(context.Background()); err != nil {
		return err
	}

	DB_POOL = db

	return nil
}
