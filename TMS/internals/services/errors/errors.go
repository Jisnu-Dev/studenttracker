package errors

import "errors"

var (
	// Domain Errors
	ErrInvalidToken          = errors.New("invalid or expired token")
	ErrSigningMethodMismatch = errors.New("unexpected signing method")

	// Fallback Function Execution Errors
	ErrGenerateTokenFailed = errors.New("unable to generate token")
	ErrValidateTokenFailed = errors.New("unable to validate token")
)
