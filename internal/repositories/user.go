package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	"snippetbox.isokol.dev/internal/models"
)

type (
	// UserRepositoryInterface - interface for user repository.
	UserRepositoryInterface interface {
		// Insert - insert new user to database.
		Insert(ctx context.Context, name, email, password string) error
		// Authenticate - verify user credentials.
		Authenticate(ctx context.Context, email, password string) (int, error)
		// Exists - check if user exists.
		Exists(ctx context.Context, id int) (bool, error)
	}
	// UserRepository - database repository for user model.
	UserRepository struct {
		// Database connection.
		db *sql.DB
	}
)

const (
	// SQL query for user insertion.
	userInsertQuery = `INSERT INTO users (name, email, hashed_password, created)
    VALUES(?, ?, ?, UTC_TIMESTAMP())`
	// SQL query to get user by email.
	userByEmailQuery = "SELECT id, hashed_password FROM users WHERE email = ?"
	// SQL query to check if user exists in database.
	userExistsQuery = "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	// Hash cost for password hashing.
	passwordHashCost = 12
	// MySQL error code for duplicated entries.
	mysqlDuplicatedErrorCode = 1062
)

// Insert - insert new user to database.
func (repository *UserRepository) Insert(ctx context.Context, name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		return fmt.Errorf("unable to hash user password: %w", err)
	}

	_, err = repository.db.ExecContext(ctx, userInsertQuery, name, email, string(hashedPassword))
	if err != nil {
		mySQLDuplicationError := checkMysqlDuplicationError(err)
		if mySQLDuplicationError != nil {
			return fmt.Errorf("database duplication error: %w", mySQLDuplicationError)
		}

		return fmt.Errorf("unable to create new user: %w", err)
	}

	return nil
}

// Check if provided error is MySQL duplication error.
func checkMysqlDuplicationError(err error) error {
	var mySQLError *mysql.MySQLError
	if errors.As(err, &mySQLError) {
		if mySQLError.Number == mysqlDuplicatedErrorCode && strings.Contains(mySQLError.Message, "users_uc_email") {
			return models.ErrDuplicateEmail
		}
	}

	return nil
}

// Authenticate - verify user credentials using provided
// email and password.
func (repository *UserRepository) Authenticate(ctx context.Context, email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	err := repository.db.QueryRowContext(ctx, userByEmailQuery, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		}

		return 0, fmt.Errorf("unable to query database for user by email: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}

		return 0, fmt.Errorf("unable to compare user password: %w", err)
	}

	return id, nil
}

// Exists - check if user exists in database.
func (repository *UserRepository) Exists(ctx context.Context, id int) (bool, error) {
	var exists bool

	err := repository.db.QueryRowContext(ctx, userExistsQuery, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error while checking if user exists: %w", err)
	}

	return exists, nil
}
