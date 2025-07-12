package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port    int    `env:"APP_PORT" envDefault:"4000"`
	BaseURL string `env:"APP_BASE_URL,expand" envDefault:"http://localhost:${APP_PORT}"`

	DatabaseUser string `env:"DATABASE_USER,notEmpty"`
	DatabasePass string `env:"DATABASE_PASSWORD,notEmpty" json:"-"`
	DatabaseHost string `env:"DATABASE_HOST,notEmpty"`
	DatabasePort int    `env:"DATABASE_PORT,notEmpty"`
	DatabaseDB   string `env:"DATABASE_DB,notEmpty"`
}

func NewFromEnv() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c Config) DBConn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.DatabaseUser, c.DatabasePass, c.DatabaseHost, c.DatabasePort, c.DatabaseDB,
	)
}

func (c Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}
