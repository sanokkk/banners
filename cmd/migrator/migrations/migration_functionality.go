package migrations

import (
	"banner-service/internal/config"
	"banner-service/internal/domain/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustMigrateWithConfig(cfg config.DbConfig) {
	const op = "migrator:MustMigrateWithConfig"

	db, err := gorm.Open(postgres.Open(cfg.ConnectionString+cfg.Database), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("%s:%w", op, err))
	}

	if err := db.AutoMigrate(&models.Banner{}); err != nil {
		panic(fmt.Errorf("%s:%w", op, err))
	}
}
