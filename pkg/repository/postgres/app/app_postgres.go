package AppPostgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

var ErrAppNotFound = errors.New("App not found")

func (p *AppPostgres) CreateApp(ctx context.Context, app models.App) (int, error) {
	// TODO
	return 0, nil
}

func (p *AppPostgres) GetApp(ctx context.Context, appid int) (models.App, error) {
	const loc = "App_Postgres.GetApp()"
	log := p.log.WithField("loc", loc)

	app := models.App{}

	query := fmt.Sprintf("SELECT id, name, token FROM %s WHERE id=$1;", app_table)
	err := p.db.Get(&app, query, appid)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("App not found")
			return models.App{}, ErrAppNotFound
		}
		log.Errorf("Error scanning app: %s", err)
		return models.App{}, err
	}

	return app, nil
}
