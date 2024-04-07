package app

import (
	"banner-service/internal/app/http"
	"banner-service/internal/config"
	"banner-service/internal/services/banner"
	"banner-service/internal/storage/local"
	repo "banner-service/internal/storage/postgres"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

type App struct {
	Server *http.Server
}

func New(cfg config.Config, logger *slog.Logger) (*App, error) {
	const op = "app:New"

	db, err := gorm.Open(postgres.Open(cfg.DbConfig.ConnectionString))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	repo := repo.New(db)
	cacheRepo := local.New()
	service := banner.New(repo, cacheRepo, logger)
	server := http.New(logger, cfg, service)

	return &App{server}, nil
}
