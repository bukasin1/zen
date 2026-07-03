package zencore

// load app config from environment variables
// it uses the [DefaultConfig] as a fallback
// Environment variables are case sensitive and should be uppercase
//
//	APP_NAME="Application name" (default: "zen-app")
//	APP_ENV="Application environment" (default: "development")
//	HTTP_ADDR="HTTP server address" (default: ":8080")
//	HTTP_SHUTDOWN_TIMEOUT="HTTP server shutdown timeout" (default: 10s)
//	LOG_LEVEL="Log level" (default: "debug")
//	LOG_PRETTY="Enable pretty printing" (default: false)
//	LOG_ENABLE_JSON="Enable JSON logging (default: false)
func LoadConfigFromEnv() Config {

	cfg := DefaultConfig()

	cfg.AppName = GetEnv(
		"APP_NAME",
		cfg.AppName,
	)

	cfg.Env = GetEnv(
		"APP_ENV",
		cfg.Env,
	)

	cfg.HTTP.Addr = GetEnv(
		"HTTP_ADDR",
		cfg.HTTP.Addr,
	)

	cfg.HTTP.ShutdownTimeout = GetEnvDuration(
		"HTTP_SHUTDOWN_TIMEOUT",
		cfg.HTTP.ShutdownTimeout,
	)

	cfg.Log.Level = GetEnv(
		"LOG_LEVEL",
		cfg.Log.Level,
	)

	cfg.Log.Pretty = GetEnvBool(
		"LOG_PRETTY",
		cfg.Log.Pretty,
	)

	cfg.Log.EnableJSON = GetEnvBool(
		"LOG_ENABLE_JSON",
		cfg.Log.EnableJSON,
	)

	return cfg
}
