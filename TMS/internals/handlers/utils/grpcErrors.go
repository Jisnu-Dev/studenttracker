package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrInvalidArgument returns a gRPC InvalidArgument (400 equivalent) error.
func ErrInvalidArgument(message string) error {
	return status.Error(codes.InvalidArgument, message)
}

// ErrInternal returns a gRPC Internal (500 equivalent) error wrapping the original error.
func ErrInternal(message string, err error) error {
	if err != nil {
		return status.Errorf(codes.Internal, "%s: %v", message, err)
	}
	return status.Error(codes.Internal, message)
}

// ErrUnauthenticated returns a gRPC Unauthenticated (401 equivalent) error.
func ErrUnauthenticated(message string) error {
	return status.Error(codes.Unauthenticated, message)
}
