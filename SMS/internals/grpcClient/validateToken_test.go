package grpcClient

import (
	"context"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
)

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name               string
		token              string
		mockAdminID        int64
		mockAdminEmail     string
		mockErr            mocks.GrpcOpError
		returnInvalidToken bool
		expectedResp       *ValidateTokenResponse
		expectedError      error
	}{
		{
			name:               "validate token successful",
			token:              "valid.jwt.token",
			mockAdminID:        42,
			mockAdminEmail:     "admin@example.com",
			expectedResp:       &ValidateTokenResponse{IsValid: true, AdminID: 42, AdminEmail: "admin@example.com"},
			expectedError:      nil,
		},
		{
			name:               "grpc call succeeds but token is invalid",
			token:              "malformed.or.tampered.token",
			returnInvalidToken: true,
			expectedResp:       &ValidateTokenResponse{IsValid: false},
			expectedError:      nil, // no gRPC error, just IsValid: false
		},
		{
			name:          "validate token fails due to invalid or expired token error",
			token:         "expired.jwt.token",
			mockErr:       mocks.GrpcOpInvalidToken,
			expectedError: mocks.ErrInvalidToken,
		},
		{
			name:          "validate token fails due to grpc internal failure",
			token:         "some.jwt.token",
			mockErr:       mocks.GrpcOpInternalError,
			expectedError: mocks.ErrValidateTokenInternal,
		},
		{
			name:          "validate token fails when TMS is unavailable",
			token:         "some.jwt.token",
			mockErr:       mocks.GrpcOpUnavailable,
			expectedError: mocks.ErrTMSUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mocks.MockTokenServiceClient{
				AdminID:            tt.mockAdminID,
				AdminEmail:         tt.mockAdminEmail,
				ValidateTokenError: tt.mockErr,
				ReturnInvalidToken: tt.returnInvalidToken,
			}

			tc := NewTokenClientForTest(mock)
			resp, err := tc.ValidateToken(context.Background(), tt.token)

			if tt.expectedError != nil {
				assertGrpcError(t, tt.expectedError, err)
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
			if tt.expectedResp != nil {
				if *resp != *tt.expectedResp {
					t.Errorf("expected response %+v, got %+v", tt.expectedResp, resp)
				}
			}

			if mock.LastValidateTokenReq == nil {
				t.Fatal("expected ValidateToken to be called on the client")
			}
			if mock.LastValidateTokenReq.Token != tt.token {
				t.Errorf("expected request token %q, got %q", tt.token, mock.LastValidateTokenReq.Token)
			}
		})
	}
}
