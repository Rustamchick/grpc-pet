package repository

import (
	"context"
	"grpc-pet/pkg/models"
	AppPostgres "grpc-pet/pkg/repository/postgres/app"
	AuthPostgres "grpc-pet/pkg/repository/postgres/auth"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Authentification interface {
	Login(ctx context.Context, email string) (user models.User, err error)
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
		Authentification: AuthPostgres.NewAuthPostgres(log, db),
		AppProvider:      AppPostgres.NewAppPostgres(log, db),
	}
}
