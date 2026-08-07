package grpcHandler

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
		expectedCode  codes.Code
		expectedError string
	}{
		{
			name:          "fails when request is nil",
			req:           nil,
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrRequestNil.Error(),
		},
		{
			name: "fails when admin id is invalid (zero)",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    0,
				AdminEmail: "admin@example.com",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrAdminIDInvalid.Error(),
		},
		{
			name: "fails when admin email is empty",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: " ",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrAdminEmailRequired.Error(),
		},
		{
			name: "fails when admin email format is invalid",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "adminexample.com",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrAdminEmailInvalid.Error(),
		},
		{
			name: "fails when admin email exceeds maximum length",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: longEmail,
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrAdminEmailTooLong.Error(),
		},
		{
			name: "validation passes with valid request",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "admin@example.com",
			},
			expectedCode: codes.OK,
		},
		{
			name: "validation passes with trimmed email",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    10,
				AdminEmail: "  admin@example.com  ",
			},
			expectedCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGenerateTokenReq(tt.req)

			if tt.expectedCode != codes.OK {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.expectedError)
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got %v", err)
				}
				if st.Code() != tt.expectedCode {
					t.Fatalf("expected code %v, got %v", tt.expectedCode, st.Code())
				}
				if st.Message() != tt.expectedError {
					t.Errorf("expected error message %q, got %q", tt.expectedError, st.Message())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateTokenReq(t *testing.T) {
	tests := []struct {
		name          string
		req           *tokenpb.ValidateTokenRequest
		expectedCode  codes.Code
		expectedError string
	}{
		{
			name:          "fails when request is nil",
			req:           nil,
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrRequestNil.Error(),
		},
		{
			name: "fails when token is empty",
			req: &tokenpb.ValidateTokenRequest{
				Token: " ",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrTokenRequired.Error(),
		},
		{
			name: "fails when token is malformed",
			req: &tokenpb.ValidateTokenRequest{
				Token: "header.payload",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrTokenMalformed.Error(),
		},
		{
			name: "validation passes with valid token structure",
			req: &tokenpb.ValidateTokenRequest{
				Token: "header.payload.signature",
			},
			expectedCode: codes.OK,
		},
		{
			name: "validation passes with trimmed whitespace token",
			req: &tokenpb.ValidateTokenRequest{
				Token: "  header.payload.signature  ",
			},
			expectedCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTokenReq(tt.req)

			if tt.expectedCode != codes.OK {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.expectedError)
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected gRPC status error, got %v", err)
				}
				if st.Code() != tt.expectedCode {
					t.Fatalf("expected code %v, got %v", tt.expectedCode, st.Code())
				}
				if st.Message() != tt.expectedError {
					t.Errorf("expected error message %q, got %q", tt.expectedError, st.Message())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

