package repository

import (
	"github.com/jmoiron/sqlx"
)

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// TODO: postgres initialization
func NewPostgresDB(cfg PostgresConfig) *sqlx.DB {
	db := sqlx.DB{}

	return &db
}
