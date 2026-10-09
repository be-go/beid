package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config is the root configuration, bound from config.json.
type Config struct {
	Gin      GinConfig      `json:"gin"`
	Server   ServerConfig   `json:"server"`
	Postgres PostgresConfig `json:"postgres"`
}

// GinConfig holds the gin settings.
type GinConfig struct {
	Mode string `json:"mode"`
}

// ServerConfig holds the HTTP server settings.
type ServerConfig struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

// PostgresConfig holds the PostgreSQL connection settings.
type PostgresConfig struct {
	DSN string `json:"dsn"`
}

// LoadConfig from config file
func LoadConfigFrom(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("read %s: %w", path, err)
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		err = fmt.Errorf("parse %s: %w", path, err)
		return nil, err
	}
	return &cfg, err
}
