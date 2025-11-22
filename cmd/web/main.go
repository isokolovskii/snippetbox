// Package main is the entry point of the web application.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type (
	application struct {
		logger *slog.Logger
		debug  bool
	}
	neuteredFileSystem struct {
		fs http.FileSystem
	}
)

const (
	exitCodeErr  = 1
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	slogKeyAddr  = "addr"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":4000"
	}

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./ui/static/"
	}

	debugStr := os.Getenv("DEBUG")
	debug, err := strconv.ParseBool(debugStr)
	if err != nil {
		debug = false
	}

	level := slog.LevelInfo

	if debug {
		level = slog.LevelDebug
	}

	handlerOpts := &slog.HandlerOptions{
		Level:       level,
		AddSource:   debug,
		ReplaceAttr: nil,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOpts))

	app := &application{
		logger: logger,
		debug:  debug,
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(staticDir),
		IdleTimeout:  time.Minute,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.InfoContext(context.Background(), "starting server", slogKeyAddr, addr)

	err = srv.ListenAndServe()
	logger.ErrorContext(context.Background(), err.Error())
	os.Exit(exitCodeErr)
}

// Open opens the named file.
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	file, err := nfs.fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	if !stat.IsDir() {
		return file, nil
	}

	file, err = nfs.openDirectory(path, file)
	if err != nil {
		return nil, fmt.Errorf("failed to open directory: %w", err)
	}

	return file, nil
}

func (nfs neuteredFileSystem) openDirectory(path string, file http.File) (http.File, error) {
	index := filepath.Join(path, "index.html")

	_, err := nfs.fs.Open(index)
	if err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("failed to close file: %w", closeErr)
		}

		return nil, fmt.Errorf("failed to open index.html: %w", err)
	}

	return file, nil
}
