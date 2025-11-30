package main

import (
	"log/slog"
	"os"
)

const (
	slogKeyIP     = "ip"
	slogKeyProto  = "proto"
	slogKeyMethod = "method"
	slogKeyURI    = "uri"
	slogKeyAddr   = "addr"
	slogKeyValue  = "value"
)

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
