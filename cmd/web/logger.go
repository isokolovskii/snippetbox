package main

import (
	"log/slog"
	"os"
)

const (
	// Log key for IP.
	slogKeyIP = "ip"
	// Log key for request protocol version.
	slogKeyProto = "proto"
	// Log key for request method.
	slogKeyMethod = "method"
	// Log key for request URI.
	slogKeyURI = "uri"
	// Log key for request address server listens to.
	slogKeyAddr = "addr"
)

// Create app logger with provided configuration.
func createLogger(loadedEnv *env) *slog.Logger {
	level := slog.LevelInfo
	if loadedEnv.debug {
		level = slog.LevelDebug
	}

	handlerOpts := &slog.HandlerOptions{
		Level:       level,
		AddSource:   loadedEnv.debug,
		ReplaceAttr: nil,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOpts))

	return logger
}
