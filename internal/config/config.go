package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	Postgres
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"5432"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"45s"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	User     string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Database string `env:"POSTGRES_DATABASE" env-default:"digital_wallet"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exists: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	return &cfg
}
