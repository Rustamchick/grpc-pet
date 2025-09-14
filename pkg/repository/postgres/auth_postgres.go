package postgres

import (
	"context"
	"fmt"
	"grpc-pet/pkg/models"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const users_table = "users"

type AuthPostgres struct {
	log *logrus.Logger
	db  *sqlx.DB
}

func NewAuthPostgres(log *logrus.Logger, db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		log: log,
		db:  db,
	}
}

func (p *AuthPostgres) Stop() error {
	return p.db.Close()
}

func (p *AuthPostgres) Login(ctx context.Context, email string, passHash []byte, appid int) (user models.User, err error) {

	return models.User{}, status.Errorf(codes.Unimplemented, "Postgres. Method Login is not implemented")
}

func (p *AuthPostgres) RegisterNewUser(ctx context.Context, email string, passHash []byte) (userid int64, err error) {
	const loc = "AuthPostgres.Register()"

	log := p.log.WithField("loc", loc)

	query := fmt.Sprintf("INSERT INTO %s (email, password_hash) VALUES ($1, $2) RETURNING id", users_table)
	row := p.db.QueryRow(query, email, passHash)

	if err := row.Scan(&userid); err != nil {

		// TODO: handle varios errors

		log.Errorf("Error scanning row: %s", err.Error())
		return 0, fmt.Errorf("Error scanning row: %s", err.Error())
	}

	return userid, nil
}

func (p *AuthPostgres) IsAdmin(ctx context.Context, userid int) (bool, error) {

	return false, status.Errorf(codes.Unimplemented, "Postgres. Method IsAdmin is not implemented")
}

func (p *AuthPostgres) App(ctx context.Context, appid int) (models.App, error) {

	return models.App{}, nil
}
