package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	TokenTTL time.Duration `yaml:"token_ttl" env-defeault:"12h"`
	GRPC     GrpcConfig    `yaml:"grpc"`
}

type GrpcConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

func InitConfig() Config {
	path := ConfigPath()
	if path == "" {
		panic("Empty config path")
	}

	return InitConfigByPath(path)
}

func InitConfigByPath(path string) Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("There is no config file in " + path)
	}

	cfg := new(Config)
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("Error reading config " + err.Error())
	}
	return *cfg
}

func ConfigPath() string {
	// в будущем можно сделать дополнительную возможность прописывать путь через флаги
	// CONFIG_PATH="D:/Proga/grpc-project/grpc-auth/configs/config.yaml"
	// res := os.Getenv("CONFIG_PATH")
	res := "config.yaml"
	return res
}
