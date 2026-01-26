package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadConfigFile reads a JSON configuration file
func LoadConfigFile(cfg *Config) {
	if cfg.ConfigPath == "" {
		cfg.ConfigPath = os.Getenv("CONFIG")
	}

	if cfg.ConfigPath == "" {
		return
	}

	file, err := os.Open(cfg.ConfigPath)
	if err != nil {
		fmt.Printf("cannot open config file: %v\n", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		fmt.Printf("cannot decode config file: %v\n", err)
	}
}
