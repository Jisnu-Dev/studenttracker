package grpcClient

import (
	"context"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name          string
		adminID       int64
		adminEmail    string
		mockToken     string
		mockErr       mocks.GrpcOpError
		expectedToken string
		expectedCode  codes.Code
		expectedError string
	}{
		{
			name:          "generate token successful",
			adminID:       1,
			adminEmail:    "admin@example.com",
			mockToken:     mocks.MockToken,
			expectedToken: mocks.MockToken,
			expectedCode:  codes.OK,
		},
		{
			name:          "generate token fails due to grpc internal failure",
			adminID:       1,
			adminEmail:    "admin@example.com",
			mockErr:       mocks.GrpcOpInternalError,
			expectedCode:  codes.Internal,
			expectedError: "failed to generate token",
		},
		{
			name:          "generate token fails when TMS is unavailable",
			adminID:       1,
			adminEmail:    "admin@example.com",
			mockErr:       mocks.GrpcOpUnavailable,
			expectedCode:  codes.Unavailable,
			expectedError: "TMS service unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mocks.MockTokenServiceClient{
				Token:       tt.mockToken,
				GenerateErr: tt.mockErr,
			}

			tc := NewTokenClientForTest(mock)
			token, err := tc.GenerateToken(context.Background(), tt.adminID, tt.adminEmail)

			if tt.expectedCode != codes.OK {
				if err == nil {
					t.Fatalf("expected error with code %v, got nil", tt.expectedCode)
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got %v", err)
				}
				if st.Code() != tt.expectedCode {
					t.Errorf("expected code %v, got %v", tt.expectedCode, st.Code())
				}
				if st.Message() != tt.expectedError {
					t.Errorf("expected message %q, got %q", tt.expectedError, st.Message())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if token != tt.expectedToken {
				t.Errorf("expected token %q, got %q", tt.expectedToken, token)
			}
		})
	}
}
