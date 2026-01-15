package AuthPostgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"grpc-pet/pkg/models"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	users_table = "grpc_users"
	apps_table  = "apps"
)

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

var (
	ErrUserExists    = errors.New("User already exist")
	ErrUserNotExists = errors.New("User not exist")
)

// Login return user by email
func (p *AuthPostgres) Login(ctx context.Context, email string) (user models.User, err error) {
	const loc = "AuthPostgres.Login()"

	log := p.log.WithField("loc", loc)

	query := fmt.Sprintf("SELECT id, email, password_hash FROM %s WHERE email=$1", users_table)
	err = p.db.Get(&user, query, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("there is no user with email: %s", email)
			return models.User{}, ErrUserNotExists
		}

		log.Errorf("Error scanning row: %s", err.Error())
		return models.User{}, fmt.Errorf("Error scanning row: %s", err.Error())
	}

	return user, nil
}

func (p *AuthPostgres) RegisterNewUser(ctx context.Context, email string, passHash []byte) (userid int64, err error) {
	const loc = "AuthPostgres.Register()"

	log := p.log.WithField("loc", loc)

	query := fmt.Sprintf("INSERT INTO %s (email, password_hash) VALUES ($1, $2) RETURNING id", users_table)
	row := p.db.QueryRow(query, email, passHash)

	if err := row.Scan(&userid); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				log.Warn("уникальный ключ нарушен (duplicate key)")
				return 0, ErrUserExists
			}
		}

		log.Errorf("Error scanning row: %s", err.Error())
		return 0, fmt.Errorf("Error scanning row: %s", err.Error())
	}

	return userid, nil
}

func (p *AuthPostgres) IsAdmin(ctx context.Context, userid int) (bool, error) {
	const loc = "AuthPostgres.IsAdmin())"

	log := p.log.WithField("loc", loc)

	isAdmin := false
	query := fmt.Sprintf("SELECT is_admin FROM %s WHERE id=$1", users_table)
	row := p.db.QueryRow(query, userid)

	if err := row.Scan(&isAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Errorf("Error no rows")
			return false, ErrUserNotExists
		}
		log.Errorf("Error scanning admin title")
		return isAdmin, err
	}
	log.Infof("userid: %d, is admin: %t", userid, isAdmin)
	return isAdmin, nil
}
