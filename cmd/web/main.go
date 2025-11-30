package main

import (
	"context"
	"log/slog"
	"net/http"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	"snippetbox.isokol.dev/internal/repositories"
)

type (
	application struct {
		logger         *slog.Logger
		repositories   *repositories.Repositories
		formDecoder    *form.Decoder
		sessionManager *scs.SessionManager
		templateCache  map[string]*template.Template
		debug          bool
	}
)

const (
	readTimeout         = 5 * time.Second
	writeTimeout        = 10 * time.Second
	databasePingTimeout = 20 * time.Second
	sessionLifetime     = 12 * time.Hour
)

func main() {
	loadedEnv := getEnv()

	logger := createLogger(loadedEnv)

	db, err := initDb(loadedEnv)
	if err != nil {
		logger.ErrorContext(context.Background(), err.Error())
		panic("Unable to open database connection")
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.ErrorContext(context.Background(), err.Error())
		panic("Unable to create template cache")
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = sessionLifetime

	app := &application{
		logger:         logger,
		debug:          loadedEnv.debug,
		repositories:   repositories.CreateRepositories(db),
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
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
