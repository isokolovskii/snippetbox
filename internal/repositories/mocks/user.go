package mocks

import (
	"context"

	"snippetbox.isokol.dev/internal/models"
)

type (
	UserRepository struct{}
)

func (*UserRepository) Insert(_ context.Context, _, email, _ string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (*UserRepository) Authenticate(_ context.Context, email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (*UserRepository) Exists(_ context.Context, id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
