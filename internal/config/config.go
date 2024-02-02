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
	Database
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"5432"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"45s"`
}

type Database struct {
	PostgresHost     string `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresDatabase string `env:"POSTGRES_DATABASE" env-default:"postgres"`
}

func MustLoad() Config {
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

	return cfg
}
