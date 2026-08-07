package grpcClient

import (
	"context"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		expectedCode    codes.Code
		expectedError   string
	}{
		{
			name:            "validate token successful",
			token:           "valid.jwt.token",
			mockAdminID:     42,
			mockEmail:       "admin@example.com",
			expectedIsValid: true,
			expectedAdminID: 42,
			expectedEmail:   "admin@example.com",
			expectedCode:    codes.OK,
		},
		{
			name:          "validate token fails due to invalid or expired token",
			token:         "invalid.jwt.token",
			mockErr:       mocks.GrpcOpInvalidToken,
			expectedCode:  codes.Unauthenticated,
			expectedError: "invalid or expired token",
		},
		{
			name:          "validate token fails due to grpc internal failure",
			token:         "some.jwt.token",
			mockErr:       mocks.GrpcOpInternalError,
			expectedCode:  codes.Internal,
			expectedError: "token validation service unavailable",
		},
		{
			name:          "validate token fails when TMS is unavailable",
			token:         "some.jwt.token",
			mockErr:       mocks.GrpcOpUnavailable,
			expectedCode:  codes.Unavailable,
			expectedError: "TMS service unavailable",
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
			resp, err := tc.ValidateToken(context.Background(), tt.token)

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
				if resp != nil {
					t.Errorf("expected nil response on error, got %+v", resp)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp == nil {
				t.Fatal("expected non-nil response")
			}
			if resp.IsValid != tt.expectedIsValid {
				t.Errorf("expected isValid %v, got %v", tt.expectedIsValid, resp.IsValid)
			}
			if resp.AdminID != tt.expectedAdminID {
				t.Errorf("expected adminID %d, got %d", tt.expectedAdminID, resp.AdminID)
			}
			if resp.AdminEmail != tt.expectedEmail {
				t.Errorf("expected email %q, got %q", tt.expectedEmail, resp.AdminEmail)
			}
		})
	}
}
