package main

import (
	app2 "banner-service/internal/app"
	"banner-service/internal/config"
	"banner-service/internal/storage/migrations"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"time"
)

const (
	envLocal  = "local"
	envDocker = "docker"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := config.MustLoad(os.Getenv("CONFIG"))
	logger := initLogger(cfg.Env)

	logger.Info("starting migrations forward")
	migrations.MustMigrateWithConfig(cfg.DbConfig)
	logger.Info("migrations ended successfully")

	logger.Info("starting application")

	app, err := app2.New(cfg, logger)
	if err != nil {
		panic(err)
	}

	tick := time.Tick(cfg.CacheConfig.TimeToUpdate)
	stopChan := make(chan os.Signal, 1)

	go app.Server.MustServe()

	for {
		select {
		case <-tick:
			app.Server.Service.UpdateCache()
			break
		case <-stopChan:
			logger.Info("application stopped") // it will be better with graceful shutdown :(
			os.Exit(1)
		}
	}
}

func initLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		break
	case envDocker:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		break
	}

	return logger
}
