package config

import (
	"log"
)

// Config variables
type Config struct {
	RunAddr    string `env:"SERVER_ADDRESS"`
	ServerAddr string `env:"BASE_URL"`
	URLData    map[string]string
}

// NewConfig create Config
func NewConfig() *Config {
	var cfg = Config{
		RunAddr:    "",
		ServerAddr: "",
		URLData:    make(map[string]string),
	}
	LoadEnv(&cfg)
	ParseFlags(&cfg, true)
	log.Println(cfg)
	return &cfg
}
