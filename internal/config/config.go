package config

import (
	"time"
)

// Config struct holds the configuration values
type Config struct {
	Timeout       time.Duration
	RetryAttempts int
}

// LoadConfig loads the config from environment variables
func LoadConfig() *Config {
	return &Config{
		Timeout:       time.Second * 10,
		RetryAttempts: 3,
	}
}
