package Postgres

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
)

type PostgresConfig struct {
	URL      string `yaml:"DATABASE_URL"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-required:"true"`
	Username string `yaml:"username" env-default:"postgres"`
	DBName   string `yaml:"dbname" env-default:"postgres"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
	Password string
}

func NewPostgresDB(cfg PostgresConfig) (*sqlx.DB, error) {
	dataSourceName := os.Getenv("DATABASE_URL") // если в докере
	if dataSourceName == "" {
		dataSourceName = cfg.URL // если локально
		if dataSourceName == "" {
			dataSourceName = fmt.Sprintf("host=%s port=%s password=%s user=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Password, cfg.Username, cfg.DBName, cfg.SSLMode) // если локально не сработал
		}
	}

	fmt.Printf("\ndataSourceName: %s\n", dataSourceName)

	db, err := sqlx.Open("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func InitPostgresConfig() PostgresConfig {
	path := dbConfigPath()
	if path == "" {
		panic("Empty config path")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("There is no config file in " + path)
	}

	cfg := new(PostgresConfig)
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("Error reading config " + err.Error())
	}

	cfg.Password = os.Getenv("POSTGRES_PASSWORD")

	return *cfg
}

func dbConfigPath() string {
	// path := os.Getenv("DB_CONFIG_PATH")
	// DB_CONFIG_PATH="D:/Proga/grpc-project/grpc-auth/configs/dbConfig.yaml"
	path := "dbConfig.yaml"
	return path
}
