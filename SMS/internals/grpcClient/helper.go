package grpcClient

import (
	"testing"

	"google.golang.org/grpc/status"
)

// assertGrpcError is a test helper that verifies if the actual error matches the expected gRPC status error.
func assertGrpcError(t *testing.T, expected, actual error) {
	t.Helper()
	if actual == nil {
		t.Fatalf("expected error %v, got nil", expected)
	}
	st, ok := status.FromError(actual)
	if !ok {
		t.Fatalf("expected gRPC status error, got %v", actual)
	}
	expectedSt, _ := status.FromError(expected)
	if st.Code() != expectedSt.Code() {
		t.Errorf("expected code %v, got %v", expectedSt.Code(), st.Code())
	}
	if st.Message() != expectedSt.Message() {
		t.Errorf("expected message %q, got %q", expectedSt.Message(), st.Message())
	}
}
