package app

import (
	grpcapp "grpc-pet/pkg/app/grpc"
	"grpc-pet/pkg/repository"
	"grpc-pet/pkg/repository/postgres"
	Service "grpc-pet/pkg/service"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(log *logrus.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// init storage
	postgresCfg := postgres.InitPostgresConfig() // TODO POSTGRES CONFIG
	db, err := postgres.NewPostgresDB(postgresCfg)
	if err != nil {
		panic(err)
	}
	log.Warn("When u make me stop?")

	repos := repository.NewRepository(log, db)

	// init auth service
	service := Service.NewService(log, repos.Authentification, repos.AppProvider, tokenTTL)

	// init app
	grpcApp := grpcapp.New(log, grpcPort, service)

	return &App{
		GRPCApp: grpcApp,
	}
}
