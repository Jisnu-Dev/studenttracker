package utils_test

import (
	"strings"
	"testing"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/Jisnu-Dev/TMS/internals/handlers/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateGenerateTokenReq(t *testing.T) {
	longEmail := strings.Repeat("a", 250) + "@example.com"

	tests := []struct {
		name        string
		req         *tokenpb.GenerateTokenRequest
		expectError bool
		wantCode    codes.Code
	}{
		{
			name:        "Nil request",
			req:         nil,
			expectError: true,
			wantCode:    codes.InvalidArgument,
		},
		{
			name: "Invalid admin id",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    0,
				AdminEmail: "admin@example.com",
			},
			expectError: true,
			wantCode:    codes.InvalidArgument,
		},
		{
			name: "Empty admin email",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: " ",
			},
			expectError: true,
			wantCode:    codes.InvalidArgument,
		},
		{
			name: "Invalid email format",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "adminexample.com",
			},
			expectError: true,
			wantCode:    codes.InvalidArgument,
		},
		{
			name: "Email exceeds maximum length",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: longEmail,
			},
			expectError: true,
			wantCode:    codes.InvalidArgument,
		},
		{
			name: "Valid request",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "admin@example.com",
			},
			expectError: false,
		},
		{
			name: "Valid request with trimmed email",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    10,
				AdminEmail: "  admin@example.com  ",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateGenerateTokenReq(tt.req)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected an error")
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got %v", err)
				}
				if st.Code() != tt.wantCode {
					t.Fatalf("expected %v, got %v", tt.wantCode, st.Code())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateValidateTokenReq(t *testing.T) {
	tests := []struct {
		name        string
		req         *tokenpb.ValidateTokenRequest
		expectError bool
		wantCode    codes.Code
	}{
		{
			name:        "Nil request",
			req:         nil,
			expectError: true,
			wantCode:    codes.InvalidArgument,
		},
		{
			name: "Empty token",
			req: &tokenpb.ValidateTokenRequest{
				Token: " ",
			},
			expectError: true,
			wantCode:    codes.InvalidArgument,
		},
		{
			name: "Malformed token",
			req: &tokenpb.ValidateTokenRequest{
				Token: "header.payload",
			},
			expectError: true,
			wantCode:    codes.InvalidArgument,
		},
		{
			name: "Valid token structure",
			req: &tokenpb.ValidateTokenRequest{
				Token: "header.payload.signature",
			},
			expectError: false,
		},
		{
			name: "Valid token with trimmed whitespace",
			req: &tokenpb.ValidateTokenRequest{
				Token: "  header.payload.signature  ",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateValidateTokenReq(tt.req)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected an error")
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got %v", err)
				}
				if st.Code() != tt.wantCode {
					t.Fatalf("expected %v, got %v", tt.wantCode, st.Code())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
