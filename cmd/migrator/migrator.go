package main

import (
	"banner-service/cmd/migrator/migrations"
	"banner-service/internal/config"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	const op = "migrator:main"
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := config.MustLoad(os.Getenv("CONFIG"))
	migrations.MustMigrateWithConfig(cfg.DbConfig)

	fmt.Println("migration ended successfully")
}
