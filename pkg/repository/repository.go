package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Authentification interface {
	Login(ctx context.Context, email, passHash string, appid int) (token string, err error)
	Register(ctx context.Context, email, passHash string) (user_id int64, err error)
	IsAdmin(ctx context.Context, userid int) (bool, error)
}

type Repository struct {
	Authentification
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authentification: NewAuthPostgres(db),
	}
}
