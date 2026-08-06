package grpcClient

import "errors"

var (
	// Fallback gRPC Operation Errors
	ErrGenerateTokenFailed = errors.New("unable to generate token")
	ErrValidateTokenFailed = errors.New("unable to validate token")
)
