package handlers

import "errors"

var (
	// Fallback Function Execution Errors
	ErrGenerateTokenFailed = errors.New("failed to sign jwt token")
	ErrValidateTokenFailed = errors.New("unable to validate token")
)
