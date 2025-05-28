package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	ErrEmptyConfig = errors.New("config is empty")
)

type (
	HTTPConfig struct {
		Port               string        `env:"HTTP_PORT"`
		ReadTimeout        time.Duration `env:"READ_TIMEOUT"`
		WriteTimeout       time.Duration `env:"WRITE_TIMEOUT"`
		MaxHeaderMegabytes int           `env:"MAX_HEADER_MBYTES"`
	}

	PostgresConfig struct {
		Host     string `env:"POSTGRES_HOST"`
		Port     int    `env:"POSTGRES_PORT"`
		Name     string `env:"POSTGRES_DB"`
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		SSLMode  string `env:"POSTGRES_SSL"`
	}

	RateLimiterConfig struct {
		Limit int `env:"LIMIT"`
	}

	Config struct {
		HTTP      HTTPConfig
		RateLimit RateLimiterConfig
		Postgres  PostgresConfig
	}
)

func New() (*Config, error) {
	var (
		err error
		cfg Config
	)

	err = cleanenv.ReadConfig(".env", &cfg)

	if cfg == (Config{}) {
		return nil, ErrEmptyConfig
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}
