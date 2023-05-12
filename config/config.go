package config

import (
	"bytes"
	"encoding/json"
	"os"
)

type Config struct {
	Database struct {
		Driver string `json:"driver"`
		DSN    string `json:"dsn"`
	} `json:"database"`
	App struct {
		Port    string `json:"port"`
		AppMode string `json:"app_mode"`
	} `json:"app"`
}

const _configPath = "./config/config.json"

func NewConfig() (*Config, error) {
	var cfg *Config
	data, err := os.ReadFile(_configPath)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(bytes.NewBuffer(data)).Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
