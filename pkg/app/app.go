package app

import (
	grpcapp "grpc-pet/pkg/app/grpc"
	"grpc-pet/pkg/repository"
	Postgres "grpc-pet/pkg/repository/postgres"
	Service "grpc-pet/pkg/service"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(log *logrus.Logger, grpcPort int, tokenTTL time.Duration) *App {
	// init storage
	postgresCfg := Postgres.InitPostgresConfig()
	db, err := Postgres.NewPostgresDB(postgresCfg)
	if err != nil {
		panic(err)
	}

	repos := repository.NewRepository(log, db)

	// init auth service
	service := Service.NewService(log, repos.Authentification, repos.AppProvider, tokenTTL)

	// init app
	grpcApp := grpcapp.New(log, grpcPort, service)

	return &App{
		GRPCApp: grpcApp,
	}
}
