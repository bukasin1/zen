package framework

import "time"

type Config struct {
	AppName string
	Env     string

	HTTP HTTPConfig
	Log  LogConfig
}

type HTTPConfig struct {
	Addr            string
	ShutdownTimeout time.Duration
}

type LogConfig struct {
	Level      string
	Pretty     bool
	EnableJSON bool
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
