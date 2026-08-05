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
		isClientError   bool
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
			token:           "custom.jwt.token",
			mockAdminID:     42,
			mockEmail:       "superadmin@school.org",
			expectedIsValid: true,
			expectedAdminID: 42,
			expectedEmail:   "superadmin@school.org",
		},
		{
			name:            "success - validly formatted token rejected by server returns false with zero values",
			token:           "invalid.jwt.token",
			mockErr:         mocks.OpInvalidToken,
			expectedIsValid: false,
			expectedAdminID: 0,
			expectedEmail:   "",
		},
		{
			name:            "success - revoked token rejected by server returns false with zero values",
			token:           "revoked.jwt.token",
			mockErr:         mocks.OpNotFound,
			expectedIsValid: false,
			expectedAdminID: 0,
			expectedEmail:   "",
		},
		{
			name:          "error - empty token string fails client validation and sends no request",
			token:         "",
			expectedError: "token is required",
			isClientError: true,
		},
		{
			name:          "error - whitespace token string fails client validation and sends no request",
			token:         "   ",
			expectedError: "token is required",
			isClientError: true,
		},
		{
			name:          "error - internal whitespace token fails client validation and sends no request",
			token:         "header.pay load.signature",
			expectedError: "token cannot contain whitespace",
			isClientError: true,
		},
		{
			name:          "error - token without dots fails client validation and sends no request",
			token:         "malformedtokenwithoutdots",
			expectedError: "malformed token structure",
			isClientError: true,
		},
		{
			name:          "error - token with wrong segment count fails client validation and sends no request",
			token:         "header.payload",
			expectedError: "malformed token structure",
			isClientError: true,
		},
		{
			name:          "error - token with whitespace and wrong segment count aggregates errors",
			token:         "foo bar",
			expectedError: "token cannot contain whitespace; malformed token structure",
			isClientError: true,
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

			if tt.isClientError {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.expectedError)
				}
				if err.Error() != tt.expectedError {
					t.Errorf("expected error %q, got %q", tt.expectedError, err.Error())
				}
				if mock.CapturedToken != "" {
					t.Errorf("expected no request sent on client error, got captured token %q", mock.CapturedToken)
				}
				return
			}

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
