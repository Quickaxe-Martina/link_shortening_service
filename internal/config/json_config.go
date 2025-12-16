package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// LoadConfigFile reads a JSON configuration file
func LoadConfigFile(cfg *Config) {
	var configPath string
	flag.StringVar(&configPath, "c", "", "config file path")
	flag.StringVar(&configPath, "config", "", "config file path")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG")
	}

	if configPath == "" {
		return
	}

	file, err := os.Open(configPath)
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
