package main

import (
	config "banner-service/internal/config"
	"banner-service/internal/domain/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	const op = "migrator:main"

	cfg := config.MustLoad().DbConfig

	db, err := gorm.Open(postgres.Open(cfg.ConnectionString), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("%s:%w", op, err))
	}

	MustMigrate(db)
	fmt.Println("migration ended successfully")
}

func MustMigrate(db *gorm.DB) {
	const op = "migrator:MustMigrate"
	if err := db.AutoMigrate(&models.Banner{}); err != nil {
		panic(fmt.Errorf("%s:%w", op, err))
	}
}
