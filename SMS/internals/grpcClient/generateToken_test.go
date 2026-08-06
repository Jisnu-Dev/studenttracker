package grpcClient

import (
	"context"
	"errors"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name          string
		adminID       int64
		adminEmail    string
		mockToken     string
		mockErr       mocks.GrpcOpError
		expectedToken string
		expectedError string
		sentinelError error
	}{
		{
			name:          "success - returns default token from grpc response",
			adminID:       1,
			adminEmail:    "admin@example.com",
			expectedToken: "mock.jwt.token",
		},
		{
			name:          "success - returns custom configured token",
			adminID:       42,
			adminEmail:    "admin.custom@example.com",
			mockToken:     "custom.signed.jwt.token",
			expectedToken: "custom.signed.jwt.token",
		},
		{
			name:          "success - handles email with special characters and plus tag",
			adminID:       5,
			adminEmail:    "admin+tag_123@sub.domain.org",
			expectedToken: "mock.jwt.token",
		},
		{
			name:          "error - grpc call fails returns sentinel error",
			adminID:       1,
			adminEmail:    "admin@example.com",
			mockErr:       mocks.GrpcOpInternalError,
			sentinelError: ErrGenerateTokenFailed,
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



			if mock.CapturedAdminID != tt.adminID {
				t.Errorf("expected captured adminID %d, got %d", tt.adminID, mock.CapturedAdminID)
			}
			if mock.CapturedAdminEmail != tt.adminEmail {
				t.Errorf("expected captured adminEmail %q, got %q", tt.adminEmail, mock.CapturedAdminEmail)
			}

			if tt.sentinelError != nil {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.sentinelError)
				}
				if !errors.Is(err, tt.sentinelError) {
					t.Errorf("expected error %q, got %q", tt.sentinelError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if token != tt.expectedToken {
				t.Errorf("expected token %q, got %q", tt.expectedToken, token)
			}
		})
	}
}
