// Package config provides a helper for loading environment variables.
package config

import (
	"github.com/battlej07/goenv"
	"github.com/joho/godotenv"
)

// Config represents configuration options for goshort.
type Config struct {
	Port        string
	BaseAddress string
}

// LoadConfig loads the configuration from the environment.
//
// It attempts to load a .env file for local development.
// Missing .env files are ignored.
// Required environment variables must be present.
func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	port, err := goenv.TryGetEnv("PORT")
	if err != nil {
		return nil, err
	}

	baseAddr, err := goenv.TryGetEnv("BASE_ADDRESS")
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:        port,
		BaseAddress: baseAddr,
	}, nil
}
