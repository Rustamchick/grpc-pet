package postgres

import (
	"context"
	"grpc-pet/pkg/models"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const app_table = "apps"

type AppPostgres struct {
	log *logrus.Logger
	db  *sqlx.DB
}

func NewAppPostgres(log *logrus.Logger, db *sqlx.DB) *AppPostgres {
	return &AppPostgres{
		log: log,
		db:  db,
	}
}

func (p *AppPostgres) CreateApp(ctx context.Context, app models.App) (int, error) {

	return 0, nil
}

func (p *AppPostgres) GetApp(ctx context.Context, appid int) (models.App, error) {

	return models.App{}, nil
}
