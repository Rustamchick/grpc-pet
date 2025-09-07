package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthPostgres struct {
	repos *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		repos: db,
	}
}

func (p *AuthPostgres) Login(ctx context.Context, email, passHash string, appid int) (token string, err error) {
	return "Postgres return token", status.Errorf(codes.Unimplemented, "Postgres. Method Login is not implemented")
}

func (p *AuthPostgres) Register(ctx context.Context, email, passHash string) (user_id int64, err error) {
	return -555, status.Errorf(codes.Unimplemented, "Postgres. Method Register is not implemented")
}

func (p *AuthPostgres) IsAdmin(ctx context.Context, userid int) (bool, error) {
	return false, status.Errorf(codes.Unimplemented, "Postgres. Method IsAdmin is not implemented")
}
