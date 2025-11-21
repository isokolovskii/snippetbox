package main

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type application struct {
	logger *slog.Logger
	debug  bool
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file, using defaults or existing environment variables")
	}

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

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level, AddSource: debug}))

	app := &application{
		logger: logger,
		debug:  debug,
	}

	logger.Info("starting server", "addr", addr)

	err = http.ListenAndServe(addr, app.routes(staticDir))
	logger.Error(err.Error())
	os.Exit(1)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)

	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}
