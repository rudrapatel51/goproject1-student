package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Address string `yaml:"address"`
}

type Postgres struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"db_name" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env-default:"disable"`
	MaxConns int32  `yaml:"max_conns" env-default:"10"`
}

type Config struct {
	Env        string `yaml:"env"`
	HTTPServer `yaml:"http_server"`
	Postgres   `yaml:"postgres"`
}

func MustLoad() *Config {
	var configPath string

	// configPath = "config/local.yaml" // path of local config file
	configPath = os.Getenv("CONFIG_PATH") // path of config file from environment variable

	if configPath == "" {
		flags := flag.String("config", "", "path to configuration file")
		flag.Parse()
		if *flags == "" {
			panic("config path is not set")
		}
		configPath = *flags

		if configPath == "" {
			panic("config path is not set")
		}
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
		panic("config file does not exist")
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err.Error())
	}

	return &cfg
}
