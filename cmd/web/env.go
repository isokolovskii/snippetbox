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
	env struct {
		addr          string
		staticDir     string
		dbDsn         string
		dbName        string
		migrationsDir string
		tlsKeyPayh    string
		tlsCertPath   string
		debug         bool
	}
)

func getEnv() *env {
	err := godotenv.Load()
	if err != nil {
		// Do nothing - try to read from env or set defaults.
		slog.Default().InfoContext(context.Background(), "no .env file, will try to get from system env or defaults")
	}

	debugStr := readEnvOrDefault("DEBUG", "false")
	debug, err := strconv.ParseBool(debugStr)
	if err != nil {
		slog.Default().WarnContext(
			context.Background(),
			"invalid DEBUG env, falling back to false",
			slogKeyValue, debugStr,
		)
		debug = false
	}

	return &env{
		addr:          readEnvOrDefault("ADDR", ":4000"),
		staticDir:     readEnvOrDefault("STATIC_DIR", "./ui/static"),
		debug:         debug,
		dbDsn:         readEnvOrDefault("DB_DSN", ""),
		migrationsDir: readEnvOrDefault("MIGRATIONS_DIR", "migrations"),
		dbName:        readEnvOrDefault("DB_NAME", "snippetbox"),
		tlsKeyPayh:    readEnvOrDefault("TLS_KEY_PATH", ""),
		tlsCertPath:   readEnvOrDefault("TLS_CERT_PATH", ""),
	}
}

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
