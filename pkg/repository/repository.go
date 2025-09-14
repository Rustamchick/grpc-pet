package repository

import (
	"context"
	"errors"
	"grpc-pet/pkg/models"
	"grpc-pet/pkg/repository/postgres"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Authentification interface {
	Login(ctx context.Context, email string, passHash []byte, appid int) (user models.User, err error)
	RegisterNewUser(ctx context.Context, email string, passHash []byte) (userid int64, err error)
	IsAdmin(ctx context.Context, userid int) (bool, error)
}

type AppProvider interface {
	CreateApp(ctx context.Context, app models.App) (int, error)
	GetApp(ctx context.Context, appid int) (models.App, error)
	// update, delete
}

type Repository struct {
	Authentification
	AppProvider
}

func NewRepository(log *logrus.Logger, db *sqlx.DB) *Repository {
	return &Repository{
		Authentification: postgres.NewAuthPostgres(log, db),
		AppProvider:      postgres.NewAppPostgres(log, db),
	}
}

var (
	ErrUserExists    = errors.New("User already exist")
	ErrUserNotExists = errors.New("User doesn't exist ")
	ErrAppNotFound   = errors.New("App not found")
)
