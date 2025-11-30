package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/go-sql-driver/mysql"

	"snippetbox.isokol.dev/migrations"
)

const (
	// Timeout for database ping.
	databasePingTimeout = 20 * time.Second
)

// Initialize database connection and run migrations.
func initDb(loadedEnv *env) (*sql.DB, error) {
	dbDsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		loadedEnv.dbUser,
		loadedEnv.dbPass,
		loadedEnv.dbHost,
		loadedEnv.dbPort,
		loadedEnv.dbName,
	)
	db, err := openDb(dbDsn)
	if err != nil {
		return nil, fmt.Errorf("unable to open connection to database: %w", err)
	}

	err = runMigrations(loadedEnv)
	if err != nil {
		defer db.Close()

		return nil, fmt.Errorf("unable to run migrations: %w", err)
	}

	return db, nil
}

// Open database connection and verify it.
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
			return nil, fmt.Errorf(
				"failed to close database connection: %w, after failure to verify database connection: %w",
				closeErr,
				err,
			)
		}

		return nil, fmt.Errorf("failed to verify connection to database: %w", err)
	}

	return db, nil
}

// Run database migrations using separate db connection.
func runMigrations(loadedEnv *env) error {
	dbDsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		loadedEnv.dbUser,
		loadedEnv.dbPass,
		loadedEnv.dbHost,
		loadedEnv.dbPort,
		loadedEnv.dbName,
	)
	db, err := openDb(dbDsn)
	if err != nil {
		return fmt.Errorf("error opening database connection for migrations: %w", err)
	}
	defer db.Close()

	databaseDriver, err := mysql.WithInstance(db, &mysql.Config{
		DatabaseName: loadedEnv.dbName,
	})
	if err != nil {
		return fmt.Errorf("error creating migration driver: %w", err)
	}

	iofsDriver, err := iofs.New(migrations.Files, ".")
	if err != nil {
		return fmt.Errorf("failed to create iofs source driver: %w", err)
	}
	defer iofsDriver.Close()

	instance, err := migrate.NewWithInstance("iofs", iofsDriver, loadedEnv.dbName, databaseDriver)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %w", err)
	}
	defer instance.Close()

	err = instance.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}
