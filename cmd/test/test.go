package main

import (
	grpcapp "grpc-pet/pkg/app/grpc"
	"grpc-pet/pkg/config"
	"grpc-pet/pkg/repository"
	Service "grpc-pet/pkg/service"
	"time"

	grpcpetv1 "github.com/Rustamchick/protobuff/gen/go/pet"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	log.SetLevel(logrus.DebugLevel)

	if err := godotenv.Load(); err != nil {
		log.Errorf("error loading env vars: %s", err)
	}

	cfg := config.InitConfig()

	postgresCfg := repository.PostgresConfig{}

	db := repository.NewPostgresDB(postgresCfg)
	repos := repository.NewRepository(db)
	authService := Service.NewService(repos)
	app := grpcapp.New(log, cfg.GRPC.Port, authService)
	app.MustRun()

	go func(app *grpcapp.App) {
		time.Sleep(time.Second * 3)
		app.Stop()
	}(&grpcapp.App{})

	TestAll(cfg, log)
}

func TestAll(cfg config.Config, Log *logrus.Logger) {
	const op = "test"
	log := Log.WithField("op", op)
	cfg_t := config.Config{
		Env:          "local",
		Storage_path: "storage_path_is_in_development",
		TokenTTL:     time.Hour * 12,
		GRPC: config.GrpcConfig{
			Port:    9090,
			Timeout: time.Hour,
		},
	}

	if cfg == cfg_t {
		log.Info("config \033[32mCONNECTED\033[0m")
	}

	req := grpcpetv1.RegisterRequest{}

	// if proto.Equal(req, req2) {}

	if req.Email == "" { // в будущем подумаю о нормальной тестировке этого пакета,
		log.Infof("Protobuff \033[32mCONNECTED\033[0m") //  но пока что по сути только вывожу красивую надпись
	}

}
