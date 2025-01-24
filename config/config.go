package config

import (
	"log"

	"github.com/caarlos0/env"
)

type Mysql struct {
	DB_HOST     string `env:"DB_HOST" envDefault:"localhost"`
	DB_USER     string `env:"DB_USER" envDefault:"vishal"`
	DB_PASSWORD string `env:"DB_PASSWORD" envDefault:"vishal"`
	DB_NAME     string `env:"DB_NAME" envDefault:"movies"`
	DB_PORT     string `env:"DB_PORT" envDefault:"3306"`
}

type Config struct{
	Mysql Mysql
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg.Mysql); err != nil {
		log.Printf("Failed to load config: %v", err)
		return nil, err
	}
	return &cfg, nil
}
