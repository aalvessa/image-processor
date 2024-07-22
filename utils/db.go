package utils

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "postgresql://imguser:f69Ac9165787@postgres:5432/imgdb?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}
