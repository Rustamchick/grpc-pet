package main

import (
	"grpc-pet/pkg/app"
	"grpc-pet/pkg/config"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	log.SetLevel(logrus.DebugLevel)

	if err := godotenv.Load(); err != nil {
		log.Errorf("error loading .env vars: %s", err)
	}

	cfg := config.InitConfig()

	application := app.New(log, cfg.GRPC.Port, cfg.Storage_path, cfg.TokenTTL)

	go application.GRPCApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	application.GRPCApp.Stop()
}
