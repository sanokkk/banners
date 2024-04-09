package config

import (
	"errors"
	"fmt"
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
	Database         string `yaml:"database"`
}

type CacheConfig struct {
	TimeToUpdate time.Duration `yaml:"time_to_update" env-default:"10ms"`
}

type ApiConfig struct {
	Port    int           `yaml:"port" env-default:"80"`
	Timeout time.Duration `yaml:"timeout" env-default:"1s"`
}

func MustLoad(path string) Config {
	const op = "config:MustLoad"

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	var config Config
	if err := yaml.Unmarshal(fileBytes, &config); err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	return config
}
