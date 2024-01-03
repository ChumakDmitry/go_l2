package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	HttpConfig `yaml:"http_config"`
}

type HttpConfig struct {
	host string `yaml:"host" env-default:"localhost"`
	port int    `yaml:"port" env-default:"8080"`
}

func ReadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config filedoes not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	fmt.Println(cfg)

	return &cfg
}
