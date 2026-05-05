package framework

import "time"

type Config struct {
	AppName string
	Env     string

	HTTP HTTPConfig
	Log  LogConfig
}

type HTTPConfig struct {
	Addr            string        // address to bind to (default: :8080)
	ShutdownTimeout time.Duration // seconds to wait for server shutdown (default: 10s)
}

type LogConfig struct {
	Level      string // log level (default: debug)
	Pretty     bool   // pretty print logs (default: false)
	EnableJSON bool   // enable JSON logging (default: false)
}

func DefaultConfig() Config {
	return Config{
		AppName: "zen-app",
		Env:     "development",
		HTTP: HTTPConfig{
			Addr:            ":8080",
			ShutdownTimeout: 10 * time.Second,
		},
		Log: LogConfig{
			Level:      "debug",
			Pretty:     false,
			EnableJSON: false,
		},
	}
}
