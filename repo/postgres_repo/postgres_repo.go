package postgres_repo

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(dbConn string) (*PostgresRepo, error) {
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepo{DB: db}, nil
}

// storer implementaion here
