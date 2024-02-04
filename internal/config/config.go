package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	SecretToket string `env:"SECRET_TOKEN" env-required:"true"`
	Env         string `yaml:"env" env-default:"local"`
	HTTPServer         //`yaml:"http_server"`
	Database
}

type HTTPServer struct {
	Host        string        `env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port        string        `env:"SERVER_PORT" env-default:"8080"`
	Timeout     time.Duration `env:"SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"SERVER_IDLETIMEOUT" env-default:"45s"`
}

type Database struct {
	PostgresHost     string `env:"POSTGRES_HOST" env-default:"0.0.0.0"`
	PostgresPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresDatabase string `env:"POSTGRES_DATABASE" env-default:"postgres"`
}

func MustLoad() Config {
	// configPath := os.Getenv("CONFIG_PATH")
	// if configPath == "" {
	// 	log.Fatal("CONFIG_PATH is not set")
	// }

	// if _, err := os.Stat(configPath); os.IsNotExist(err) {
	// 	log.Fatalf("config file doesn't exists: %s", configPath)
	// }

	var cfg Config
	// if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
	// 	log.Fatalf("can't read config: %s", err)
	// }

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	return cfg
}
