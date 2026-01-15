package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var method string

	flag.StringVar(&method, "method", "", "up or down")
	flag.Parse()

	migrationsPath := "./migrations"

	m, err := migrate.New("file://"+migrationsPath, "postgres://postgres:12345@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	switch method {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")
				return
			}
			panic(err)
		}
		fmt.Printf("%s migrations applied ✅", method)
	case "down":
		if err := m.Down(); err != nil {
			panic(err)
		}
		fmt.Printf("%s migrations applied ✅", method)
	}
}
