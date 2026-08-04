package grpcClient

import (
	"context"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name          string
		adminID       int64
		adminEmail    string
		mockToken     string
		mockErr       mocks.MockOpError
		expectedToken string
		expectedError string
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
			name:          "success - handles empty email boundary",
			adminID:       10,
			adminEmail:    "",
			expectedToken: "mock.jwt.token",
		},
		{
			name:          "success - handles zero admin ID boundary",
			adminID:       0,
			adminEmail:    "zero@example.com",
			expectedToken: "mock.jwt.token",
		},
		{
			name:          "success - handles negative admin ID boundary",
			adminID:       -1,
			adminEmail:    "neg@example.com",
			expectedToken: "mock.jwt.token",
		},
		{
			name:          "success - handles email with special characters and plus tag",
			adminID:       5,
			adminEmail:    "admin+tag_123@sub.domain.org",
			expectedToken: "mock.jwt.token",
		},
		{
			name:          "error - grpc call fails returns wrapped error",
			adminID:       1,
			adminEmail:    "admin@example.com",
			mockErr:       mocks.OpInternalError,
			expectedError: "failed to generate token: rpc error: failed to generate token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mocks.MockTokenServiceClient{
				Token:              tt.mockToken,
				GenerateTokenError: tt.mockErr,
			}

			tc := newTokenClientForTest(mock)
			token, err := tc.GenerateToken(context.Background(), tt.adminID, tt.adminEmail)

			if mock.CapturedAdminID != tt.adminID {
				t.Errorf("expected captured adminID %d, got %d", tt.adminID, mock.CapturedAdminID)
			}
			if mock.CapturedAdminEmail != tt.adminEmail {
				t.Errorf("expected captured adminEmail %q, got %q", tt.adminEmail, mock.CapturedAdminEmail)
			}

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.expectedError)
				}
				if err.Error() != tt.expectedError {
					t.Errorf("expected error %q, got %q", tt.expectedError, err.Error())
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
