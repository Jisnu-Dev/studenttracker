package handlers

import (
	"strings"
	"testing"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateGenerateTokenReq(t *testing.T) {
	longEmail := strings.Repeat("a", 250) + "@example.com"

	tests := []struct {
		name          string
		req           *tokenpb.GenerateTokenRequest
		wantCode      codes.Code
		expectedError string
	}{
		{
			name:          "Nil request",
			req:           nil,
			wantCode:      codes.InvalidArgument,
			expectedError: "request cannot be nil",
		},
		{
			name: "Invalid admin id",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    0,
				AdminEmail: "admin@example.com",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_id is required and must be greater than 0",
		},
		{
			name: "Empty admin email",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: " ",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_email is required",
		},
		{
			name: "Invalid email format",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "adminexample.com",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_email is invalid",
		},
		{
			name: "Email exceeds maximum length",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: longEmail,
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_email exceeds maximum length",
		},
		{
			name: "Valid request",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "admin@example.com",
			},
		},
		{
			name: "Valid request with trimmed email",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    10,
				AdminEmail: "  admin@example.com  ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGenerateTokenReq(tt.req)

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.expectedError)
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got %v", err)
				}
				if st.Code() != tt.wantCode {
					t.Fatalf("expected code %v, got %v", tt.wantCode, st.Code())
				}
				if st.Message() != tt.expectedError {
					t.Errorf("expected error message %q, got %q", tt.expectedError, st.Message())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateValidateTokenReq(t *testing.T) {
	tests := []struct {
		name          string
		req           *tokenpb.ValidateTokenRequest
		wantCode      codes.Code
		expectedError string
	}{
		{
			name:          "Nil request",
			req:           nil,
			wantCode:      codes.InvalidArgument,
			expectedError: "request cannot be nil",
		},
		{
			name: "Empty token",
			req: &tokenpb.ValidateTokenRequest{
				Token: " ",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "token is required",
		},
		{
			name: "Malformed token",
			req: &tokenpb.ValidateTokenRequest{
				Token: "header.payload",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "malformed token structure",
		},
		{
			name: "Valid token structure",
			req: &tokenpb.ValidateTokenRequest{
				Token: "header.payload.signature",
			},
		},
		{
			name: "Valid token with trimmed whitespace",
			req: &tokenpb.ValidateTokenRequest{
				Token: "  header.payload.signature  ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateValidateTokenReq(tt.req)

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.expectedError)
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got %v", err)
				}
				if st.Code() != tt.wantCode {
					t.Fatalf("expected code %v, got %v", tt.wantCode, st.Code())
				}
				if st.Message() != tt.expectedError {
					t.Errorf("expected error message %q, got %q", tt.expectedError, st.Message())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestResponsesAndGrpcError(t *testing.T) {
	t.Run("GrpcError helper", func(t *testing.T) {
		err := GrpcError(codes.NotFound, "not found test")
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected grpc status error")
		}
		if st.Code() != codes.NotFound || st.Message() != "not found test" {
			t.Errorf("unexpected status: %v", st)
		}
	})

	t.Run("GenerateTokenResponse helper", func(t *testing.T) {
		resp := GenerateTokenResponse("sample-token")
		if resp == nil || resp.GetToken() != "sample-token" {
			t.Errorf("unexpected response: %v", resp)
		}
	})

	t.Run("ValidateTokenResponse helper", func(t *testing.T) {
		resp := ValidateTokenResponse(true, 42, "admin@test.com")
		if resp == nil || !resp.GetIsValid() || resp.GetAdminId() != 42 || resp.GetAdminEmail() != "admin@test.com" {
			t.Errorf("unexpected response: %v", resp)
		}
	})
}
