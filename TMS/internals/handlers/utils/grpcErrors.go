package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcError creates and returns a gRPC status error from a status code and error message.
func GrpcError(code codes.Code, message string) error {
	return status.Error(code, message)
}

