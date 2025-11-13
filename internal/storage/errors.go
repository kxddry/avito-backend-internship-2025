package storage

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Errors are the errors for the storage.
var (
	// ErrNotFound is the error for not found.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is the error for already exists.
	ErrAlreadyExists = errors.New("already exists")
)

// Unique violation error code.
const (
	errUniqueViolation = "23505"
)

// IsUniqueViolation checks if the error is a unique violation error.
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == errUniqueViolation
	}
	return false
}
