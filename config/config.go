package config

import (
	"strings"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Env        string `env:"TODO_ENV" envDefault:"dev"`
	Port       int    `env:"PORT" envDefault:"80"`
	DBHost     string `env:"TODO_DB_HOST" envDefault:"localhost"`
	DBPort     int    `env:"TODO_DB_PORT" envDefault:"33306"`
	DBUser     string `env:"TODO_DB_USER" envDefault:"todo"`
	DBPassword string `env:"TODO_DB_PASSWORD" envDefault:"todo"`
	DBName     string `env:"TODO_DB_NAME" envDefault:"todo"`
	RedisHost  string `env:"TODO_REDIS_HOST" envDefault:"127.0.0.1"`
	RedisPort  int    `env:"TODO_REDIS_PORT" envDefault:"36379"`
	// CORS設定
	CORSAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
	CORSAllowedMethods string `env:"CORS_ALLOWED_METHODS" envDefault:"GET,POST,PUT,DELETE,OPTIONS"`
	CORSAllowedHeaders string `env:"CORS_ALLOWED_HEADERS" envDefault:"Accept,Authorization,Content-Type,X-CSRF-Token,X-Requested-With"`
	CORSMaxAge         int    `env:"CORS_MAX_AGE" envDefault:"86400"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// CORSOptions CORS設定オプション
type CORSOptions struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
}

// GetCORSOptions 設定からCORSオプションを取得
func (c *Config) GetCORSOptions() *CORSOptions {
	return &CORSOptions{
		AllowedOrigins: strings.Split(c.CORSAllowedOrigins, ","),
		AllowedMethods: strings.Split(c.CORSAllowedMethods, ","),
		AllowedHeaders: strings.Split(c.CORSAllowedHeaders, ","),
		MaxAge:         c.CORSMaxAge,
	}
}
