package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

var (
	ErrConfigPathEmpty = errors.New("config path from environment is empty")
)

type Config struct {
	Env         string      `yaml:"env" env-required:"true"`
	ApiConfig   ApiConfig   `yaml:"api"`
	DbConfig    DbConfig    `yaml:"db"`
	CacheConfig CacheConfig `yaml:"cache_config"`
}

type DbConfig struct {
	ConnectionString string `yaml:"connection_string" env-required:"true"`
}

type CacheConfig struct {
	TimeToUpdate time.Duration `yaml:"time_to_update" env-default:"10ms"`
}

type ApiConfig struct {
	Port    int           `yaml:"port" env-default:"80"`
	Timeout time.Duration `yaml:"timeout" env-default:"1s"`
}

func MustLoad() Config {
	const op = "config:MustLoad"
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		panic(fmt.Errorf("%s: %w", op, ErrConfigPathEmpty))
	}

	fileBytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	var config Config
	if err := yaml.Unmarshal(fileBytes, &config); err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	return config
}
