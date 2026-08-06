package grpcClient

import (
	"context"
	"errors"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
)

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		mockAdminID     int64
		mockEmail       string
		mockErr         mocks.GrpcOpError
		expectedIsValid bool
		expectedAdminID int64
		expectedEmail   string
		expectedError   string
		sentinelError   error
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
			name:          "error - invalid token rejected by server returns sentinel error",
			token:         "invalid.jwt.token",
			mockErr:       mocks.GrpcOpInvalidToken,
			sentinelError: ErrValidateTokenFailed,
		},
		{
			name:          "error - grpc internal failure returns sentinel error",
			token:         "some.jwt.token",
			mockErr:       mocks.GrpcOpInternalError,
			sentinelError: ErrValidateTokenFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mocks.MockTokenServiceClient{
				AdminID:     tt.mockAdminID,
				AdminEmail:  tt.mockEmail,
				ValidateErr: tt.mockErr,
			}

			tc := NewTokenClientForTest(mock)
			isValid, adminID, email, err := tc.ValidateToken(context.Background(), tt.token)



			if mock.CapturedToken != tt.token {
				t.Errorf("expected captured token %q, got %q", tt.token, mock.CapturedToken)
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
