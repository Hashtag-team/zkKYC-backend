package db

import (
	"database/sql"
	"zkKYC-backend/internal/app/config"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func NewDBConnection(config config.Config) *sql.DB {
	db, err := sql.Open("pgx", config.DatabaseDSN)

	if err != nil {
		panic(err)
	}

	return db
}
