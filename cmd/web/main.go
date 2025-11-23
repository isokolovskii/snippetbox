package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"

	"snippetbox.isokol.dev/internal/repositories"
)

type (
	application struct {
		logger       *slog.Logger
		repositories *repositories.Repositories
		debug        bool
	}
	neuteredFileSystem struct {
		fs http.FileSystem
	}
	env struct {
		addr      string
		staticDir string
		dbDsn     string
		debug     bool
	}
)

const (
	readTimeout         = 5 * time.Second
	writeTimeout        = 10 * time.Second
	databasePingTimeout = 20 * time.Second
	slogKeyAddr         = "addr"
	slogKeyValue        = "value"
)

func main() {
	loadedEnv := getEnv()

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

	db, err := openDb(loadedEnv.dbDsn)
	if err != nil {
		logger.ErrorContext(context.Background(), err.Error())
		panic("Unable to open database connection")
	}

	defer db.Close()

	app := &application{
		logger:       logger,
		debug:        loadedEnv.debug,
		repositories: repositories.CreateRepositories(db),
	}

	srv := &http.Server{
		Addr:         loadedEnv.addr,
		Handler:      app.routes(loadedEnv.staticDir),
		IdleTimeout:  time.Minute,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.InfoContext(context.Background(), "starting server", slogKeyAddr, loadedEnv.addr)

	err = srv.ListenAndServe()
	logger.ErrorContext(context.Background(), err.Error())
	panic("unexpected server failure")
}

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

func getEnv() *env {
	err := godotenv.Load()
	if err != nil {
		// Do nothing - try to read from env or set defaults.
		slog.Default().InfoContext(context.Background(), "no .env file, will try to get from system env or defaults")
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
	debug := false
	if debugStr != "" {
		debug, err = strconv.ParseBool(debugStr)
		if err != nil {
			slog.Default().WarnContext(
				context.Background(),
				"invalid DEBIG env, falling back to false",
				slogKeyValue, debugStr,
			)
			debug = false
		}
	}

	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		slog.Default().ErrorContext(context.Background(), "no DB_DSN provided from env")
		panic("No database connection info")
	}

	return &env{
		addr:      addr,
		staticDir: staticDir,
		debug:     debug,
		dbDsn:     dbDSN,
	}
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

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), databasePingTimeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		closeErr := db.Close()

		if closeErr != nil {
			return nil, fmt.Errorf("failed to close database connection: %w", closeErr)
		}

		return nil, fmt.Errorf("failed to verify connection to database: %w", err)
	}

	return db, nil
}
