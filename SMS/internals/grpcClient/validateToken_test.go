package grpcClient

import (
	"context"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
)

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		mockAdminID     int64
		mockEmail       string
		mockErr         mocks.MockOpError
		expectedIsValid bool
		expectedAdminID int64
		expectedEmail   string
		expectedError   string
	}{
		{
			name:            "success - valid token returns default claims",
			token:           "valid.jwt.token",
			expectedIsValid: true,
			expectedAdminID: 1,
			expectedEmail:   "admin@example.com",
		},
		{
			name:            "success - valid token returns custom claims",
			token:           "valid.custom.jwt.token",
			mockAdminID:     42,
			mockEmail:       "superadmin@school.org",
			expectedIsValid: true,
			expectedAdminID: 42,
			expectedEmail:   "superadmin@school.org",
		},
		{
			name:            "success - invalid token returns false with zero values",
			token:           "invalid.jwt.token",
			mockErr:         mocks.OpInvalidToken,
			expectedIsValid: false,
			expectedAdminID: 0,
			expectedEmail:   "",
		},
		{
			name:            "success - not found or revoked token returns false with zero values",
			token:           "revoked.jwt.token",
			mockErr:         mocks.OpNotFound,
			expectedIsValid: false,
			expectedAdminID: 0,
			expectedEmail:   "",
		},
		{
			name:            "success - empty token string returns false when marked invalid",
			token:           "",
			mockErr:         mocks.OpInvalidToken,
			expectedIsValid: false,
			expectedAdminID: 0,
			expectedEmail:   "",
		},
		{
			name:            "success - malformed token string with special characters",
			token:           "Bearer !@#$%^&*()_+-=[]{}|;':,.<>/?",
			mockErr:         mocks.OpInvalidToken,
			expectedIsValid: false,
			expectedAdminID: 0,
			expectedEmail:   "",
		},
		{
			name:          "error - grpc call fails returns wrapped error",
			token:         "some.jwt.token",
			mockErr:       mocks.OpInternalError,
			expectedError: "failed to validate token: rpc error: token validation service unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mocks.MockTokenServiceClient{
				AdminID:            tt.mockAdminID,
				AdminEmail:         tt.mockEmail,
				ValidateTokenError: tt.mockErr,
			}

			tc := newTokenClientForTest(mock)
			isValid, adminID, email, err := tc.ValidateToken(context.Background(), tt.token)

			if mock.CapturedToken != tt.token {
				t.Errorf("expected captured token %q, got %q", tt.token, mock.CapturedToken)
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
			if isValid != tt.expectedIsValid {
				t.Errorf("expected isValid %v, got %v", tt.expectedIsValid, isValid)
			}
			if adminID != tt.expectedAdminID {
				t.Errorf("expected adminID %d, got %d", tt.expectedAdminID, adminID)
			}
			if email != tt.expectedEmail {
				t.Errorf("expected email %q, got %q", tt.expectedEmail, email)
			}
		})
	}
}
