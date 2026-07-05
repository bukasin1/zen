package zencore

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func GetEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
		return fallback
	}

	return value
}

func GetEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func GetEnvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func GetEnvDuration(
	key string,
	fallback time.Duration,
) time.Duration {

	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func MustGetEnv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
		panic(
			newFrameworkPanic(
				"missing required environment variable: " + key,
			),
		)
	}

	return value
}

func MustGetEnvInt(key string) int {
	value := MustGetEnv(key)

	parsed, err := strconv.Atoi(value)
	if err != nil {
		panic(
			newFrameworkPanic(
				"invalid integer environment variable: " + key,
			),
		)
	}

	return parsed
}

func MustGetEnvBool(key string) bool {
	value := MustGetEnv(key)

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		panic(
			newFrameworkPanic(
				"invalid boolean environment variable: " + key,
			),
		)
	}

	return parsed
}

func MustGetEnvDuration(key string) time.Duration {
	value := MustGetEnv(key)

	parsed, err := time.ParseDuration(value)
	if err != nil {
		panic(
			newFrameworkPanic(
				"invalid duration environment variable: " + key,
			),
		)
	}

	return parsed
}
