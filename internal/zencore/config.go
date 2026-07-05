package zencore

import "time"

type Config struct {
	AppName string
	Env     string

	HTTP HTTPConfig
	Log  LogConfig
}

type HTTPConfig struct {
	Addr              string        // address to bind to (default: :8080)
	ReadTimeout       time.Duration // duration to wait for request to be read
	ReadHeaderTimeout time.Duration // duration to wait for request headers to be read
	WriteTimeout      time.Duration // duration to wait for response to be written
	IdleTimeout       time.Duration // duration to wait for next request
	MaxHeaderBytes    int           // maximum number of bytes to read from request headers
	MaxBodyBytes      int64         // maximum number of bytes to read from request body
	ShutdownTimeout   time.Duration // seconds to wait for server shutdown (default: 10s)
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
