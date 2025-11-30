package main

import (
	"context"
	"crypto/tls"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	"snippetbox.isokol.dev/internal/repositories"
)

type (
	// Application helpers provided as DI.
	application struct {
		// Application logger.
		logger *slog.Logger
		// Database repositories.
		repositories *repositories.Repositories
		// Form decoder instance.
		formDecoder *form.Decoder
		// Session manager instance.
		sessionManager *scs.SessionManager
		// Rendered templates cache.
		templateCache map[string]*template.Template
		// Server debig config.
		debug bool
	}
)

const (
	// Server request read timeout.
	readTimeout = 5 * time.Second
	// Server response write timeout.
	writeTimeout = 10 * time.Second
	// Session lifetime duration.
	sessionLifetime = 12 * time.Hour
)

// Server bootstrap. Creates all entities required
// for server work. Also connects to database and
// sets configuration for http server.
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
	sessionManager.Cookie.Secure = true

	app := &application{
		logger:         logger,
		debug:          loadedEnv.debug,
		repositories:   repositories.CreateRepositories(db),
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         loadedEnv.addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
	}

	logger.InfoContext(context.Background(), "starting server", slogKeyAddr, loadedEnv.addr)

	err = srv.ListenAndServeTLS(loadedEnv.tlsCertPath, loadedEnv.tlsKeyPath)
	logger.ErrorContext(context.Background(), err.Error())
	panic("unexpected server failure")
}
