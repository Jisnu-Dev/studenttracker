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
		isClientError bool
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
			name:          "error - empty email fails client validation and sends no request",
			adminID:       10,
			adminEmail:    "",
			expectedError: "admin_email: email is required",
			isClientError: true,
		},
		{
			name:          "error - whitespace email fails client validation and sends no request",
			adminID:       10,
			adminEmail:    "   ",
			expectedError: "admin_email: email cannot be empty or only whitespace",
			isClientError: true,
		},
		{
			name:          "error - invalid email format fails client validation and sends no request",
			adminID:       10,
			adminEmail:    "not-an-email",
			expectedError: "admin_email: email must contain '@'",
			isClientError: true,
		},
		{
			name:          "error - zero admin ID fails client validation and sends no request",
			adminID:       0,
			adminEmail:    "zero@example.com",
			expectedError: "admin_id is required and must be greater than 0",
			isClientError: true,
		},
		{
			name:          "error - negative admin ID fails client validation and sends no request",
			adminID:       -1,
			adminEmail:    "neg@example.com",
			expectedError: "admin_id is required and must be greater than 0",
			isClientError: true,
		},
		{
			name:          "error - both zero admin ID and empty email fail client validation and send no request",
			adminID:       0,
			adminEmail:    "",
			expectedError: "admin_id is required and must be greater than 0; admin_email: email is required",
			isClientError: true,
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

			if tt.isClientError {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.expectedError)
				}
				if err.Error() != tt.expectedError {
					t.Errorf("expected error %q, got %q", tt.expectedError, err.Error())
				}
				if mock.CapturedAdminID != 0 || mock.CapturedAdminEmail != "" {
					t.Errorf("expected no request sent on client error, got captured adminID %d, email %q", mock.CapturedAdminID, mock.CapturedAdminEmail)
				}
				return
			}

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
