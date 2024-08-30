package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/germandv/apio/internal/config"
	"github.com/germandv/apio/internal/db"
)

type MigratorConfig struct {
	PostgresConnStr string `env:"POSTGRES_CONN_STR"`
}

func main() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	cfg := MigratorConfig{}
	err := config.Parse(&cfg, envFile)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	migrator, err := db.NewDBMigrator(cfg.PostgresConnStr, os.DirFS("migrations"))
	if err != nil {
		panic(err)
	}

	v, err := migrator.Status(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Version Before: %d\n", v)

	var action string
	flag.StringVar(&action, "action", "up", "action to perform (up|down)")
	flag.Parse()

	switch action {
	case "up":
		err = migrator.Up(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println("Successfully run UP migrations")
	case "down":
		err = migrator.Down(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println("Successfully undid last migration")
	default:
		panic("action must be 'up' or 'down'")
	}

	v, err = migrator.Status(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Version After: %d\n", v)
}
