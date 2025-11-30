package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	// Environment based config.
	env struct {
		// Server address.
		addr string
		// Database host. Required value in .env or environment.
		dbHost string
		// Database user. Required value in .env or environment.
		dbUser string
		// Database user password. Required value in .env or environment.
		dbPass string
		// Database name. Required value in .env or environment.
		dbName string
		// TLS key path. Required value in .env or environment.
		tlsKeyPath string
		// TLS certificate path. Required value in .env or environment.
		tlsCertPath string
		// Run server in debug mode.
		debug bool
	}
)

// Reads env variables from .env or from system environment
// DB_HOST, DB_USER, DB_PASS, TLS_KEY_PATH and TLS_CERT_PATH are required variables
// If required variables not provided via environment this function
// will panic.
func getEnv() *env {
	err := godotenv.Load()
	if err != nil {
		// Do nothing - try to read from env or set defaults.
		slog.Default().InfoContext(context.Background(), "no .env file, will try to get from system env or defaults")
	}

	return &env{
		addr:        readEnvOrDefault("ADDR", ":4000"),
		debug:       parseEnvBool("DEBUG", "false"),
		dbHost:      readEnvOrDefault("DB_HOST", ""),
		dbUser:      readEnvOrDefault("DB_USER", ""),
		dbPass:      readEnvOrDefault("DB_PASS", ""),
		dbName:      readEnvOrDefault("DB_NAME", "snippetbox"),
		tlsKeyPath:  readEnvOrDefault("TLS_KEY_PATH", ""),
		tlsCertPath: readEnvOrDefault("TLS_CERT_PATH", ""),
	}
}

// Reads variable from environment by provided key.
// If variable not found in environment - will use
// default variable. If empty string is provided as
// default this function assumes that this variable
// is required and no sensible default may exist.
// If required variable is not provided this will be
// logged and function will panic.
func readEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		if defaultValue != "" {
			return defaultValue
		}
		slog.Default().ErrorContext(context.Background(), fmt.Sprintf("no %s provided from env", key))
		panic("required env value not provided")
	}

	return value
}

// Parse specified env variable as boolean and will panic for unprocessable values.
func parseEnvBool(key, defaultValue string) bool {
	valueStr := readEnvOrDefault(key, defaultValue)
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		panic(fmt.Sprintf("invalid %s env, should be `true` of `false`, got %s", key, valueStr))
	}

	return value
}
