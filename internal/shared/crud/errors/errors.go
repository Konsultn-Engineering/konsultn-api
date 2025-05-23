package errors

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// Common repository errors
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrDatabaseError  = errors.New("database error")
	ErrContextTimeout = errors.New("context timeout or cancelled")
)

// WrapError wraps a database error with a more descriptive message
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}

	// Handle specific cases
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%s: %w", message, ErrRecordNotFound)
	}

	// Check for context errors
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%s: %w", message, ErrContextTimeout)
	}

	return fmt.Errorf("%s: %w", message, err)
}
