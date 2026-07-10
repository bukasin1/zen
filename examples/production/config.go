package main

import (
	"time"

	"github.com/bukasin1/zen"
)

// Config contains the application's configuration.
type Config struct {
	Server struct {
		Address string
	}

	HTTP struct {
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}

	Logging struct {
		JSON bool
	}
}

// loadConfig loads configuration from environment variables.
func loadConfig() Config {
	cfg := Config{}

	cfg.Server.Address = zen.GetEnv("SERVER_ADDRESS", ":8080")

	cfg.HTTP.ReadTimeout = zen.GetEnvDuration(
		"HTTP_READ_TIMEOUT",
		5*time.Second,
	)

	cfg.HTTP.WriteTimeout = zen.GetEnvDuration(
		"HTTP_WRITE_TIMEOUT",
		10*time.Second,
	)

	cfg.HTTP.IdleTimeout = zen.GetEnvDuration(
		"HTTP_IDLE_TIMEOUT",
		60*time.Second,
	)

	cfg.Logging.JSON = zen.GetEnvBool(
		"LOG_JSON",
		false,
	)

	return cfg
}
