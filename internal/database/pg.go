package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool() (*pgxpool.Pool, error) {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgres://username:password@localhost:5432/mydb"
	}

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS details (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username TEXT NOT NULL,
			phone TEXT NOT NULL UNIQUE,
			createdAt TIMESTAMP DEFAULT NOW(),
			updatedAt TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
