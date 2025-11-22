package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env string `env:"ENV" env-default:"local"`
	HTTPConfig
	DBConfig
}
type HTTPConfig struct {
	HTTPHost string `env:"HTTP_HOST" env-default:"localhost"`
	HTTPPort int    `env:"HTTP_PORT" env-default:"8080"`
}
type DBConfig struct {
	DBName string `env:"DB_NAME"`
	DBPort int    `env:"DB_PORT"`
	DBHost string `env:"DB_HOST"`
	DBPass string `env:"DB_PASS"`
	DBUser string `env:"DB_USER"`
}

func New() (*Config, error) {
	var cfg Config
	config_path := os.Getenv("ENV_FILE")
	if config_path == "" {
		config_path = "../.env"
	}
	_ = godotenv.Load(config_path)

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName)
}
