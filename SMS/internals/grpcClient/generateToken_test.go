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
		mockErr       mocks.GrpcOpError
		expectedToken string
		expectedError error
	}{
		{
			name:          "generate token successful",
			adminID:       1,
			adminEmail:    "admin@example.com",
			expectedToken: mocks.MockToken,
			expectedError: nil,
		},
		{
			name:          "generate token fails due to grpc internal failure",
			adminID:       1,
			adminEmail:    "admin@example.com",
			mockErr:       mocks.GrpcOpInternalError,
			expectedError: mocks.ErrGenerateTokenInternal,
		},
		{
			name:          "generate token fails when TMS is unavailable",
			adminID:       1,
			adminEmail:    "admin@example.com",
			mockErr:       mocks.GrpcOpUnavailable,
			expectedError: mocks.ErrTMSUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mocks.MockTokenServiceClient{
				GenerateTokenError: tt.mockErr,
			}

			tc := NewTokenClientForTest(mock)
			token, err := tc.GenerateToken(context.Background(), tt.adminID, tt.adminEmail)

			if tt.expectedError != nil {
				assertGrpcError(t, tt.expectedError, err)
				if token != "" {
					t.Errorf("expected empty token on error, got %q", token)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if token != tt.expectedToken {
				t.Errorf("expected token %q, got %q", tt.expectedToken, token)
			}

			// Verify the request that actually went out over the wire.
			if mock.LastGenerateTokenReq == nil {
				t.Fatal("expected GenerateToken to be called on the client")
			}
			if mock.LastGenerateTokenReq.AdminID != tt.adminID {
				t.Errorf("expected request AdminID %d, got %d", tt.adminID, mock.LastGenerateTokenReq.AdminID)
			}
			if mock.LastGenerateTokenReq.AdminEmail != tt.adminEmail {
				t.Errorf("expected request AdminEmail %q, got %q", tt.adminEmail, mock.LastGenerateTokenReq.AdminEmail)
			}
		})
	}
}
