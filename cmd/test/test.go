package main

import (
	"fmt"
	"grpc-pet/pkg/config"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcpetv1 "github.com/Rustamchick/protobuff/gen/go/pet"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func main() {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	log.SetLevel(logrus.DebugLevel)

	if err := godotenv.Load(); err != nil {
		log.Errorf("error loading env vars: %s", err)
	}

	cfg := config.InitConfig()

	// application := app.New(log, cfg.GRPC.Port, cfg.Storage_path, cfg.TokenTTL)

	// app := grpcapp.New(log, cfg.GRPC.Port, authService)

	// go app.MustRun()

	TestAll(cfg, log)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
}

func PreStart() {

}

func testProto() {
	newReq := &grpcpetv1.RegisterRequest{
		Email:    "Rustam@example.com",
		Password: "12345",
	}
	fmt.Println(newReq)
	// сериализация данных
	data, err := proto.Marshal(newReq.ProtoReflect().Interface())
	if err != nil {
		log.Fatal("Marshaling error: ", err)
	}

	fmt.Println(data)
	// десериализация данных
	newReq2 := &grpcpetv1.RegisterRequest{}
	err = proto.Unmarshal(data, newReq2)
	if err != nil {
		log.Fatal("Unmarshaling error: ", err)
	}

	fmt.Println(newReq2)
}

func TestAll(cfg config.Config, Log *logrus.Logger) {
	const loc = "test"
	log := Log.WithField("loc", loc)

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
