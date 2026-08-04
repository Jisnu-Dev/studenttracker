package handlers

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// assertGRPCError fails the test unless err is a gRPC status error with the
// expected code, and whose message exactly matches expectedError.
func assertGRPCError(t *testing.T, err error, wantCode codes.Code, expectedError string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected a gRPC status error, got %v", err)
	}

	if st.Code() != wantCode {
		t.Errorf("expected code %v, got %v (message: %q)", wantCode, st.Code(), st.Message())
	}

	if st.Message() != expectedError {
		t.Errorf("expected error message %q, got %q", expectedError, st.Message())
	}
}

// assertValidToken parses resp's token with secret and verifies its claims
// match what was requested, including a ~5 hour expiry window.
func assertValidToken(t *testing.T, resp *tokenpb.GenerateTokenResponse, secret []byte, wantAdminID int64, wantAdminEmail string) {
	t.Helper()

	if resp == nil || resp.GetToken() == "" {
		t.Fatal("expected a non-empty token in the response")
	}

	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(resp.GetToken(), claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		t.Fatalf("failed to parse generated token: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("expected generated token to be valid")
	}

	if claims.AdminID != wantAdminID {
		t.Errorf("expected AdminID %d, got %d", wantAdminID, claims.AdminID)
	}
	if claims.AdminEmail != wantAdminEmail {
		t.Errorf("expected AdminEmail %q, got %q", wantAdminEmail, claims.AdminEmail)
	}

	if claims.IssuedAt == nil {
		t.Fatal("expected IssuedAt to be set")
	}
	if claims.ExpiresAt == nil {
		t.Fatal("expected ExpiresAt to be set")
	}

	const wantTTL = 5 * time.Hour
	gotTTL := claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time)
	if diff := gotTTL - wantTTL; diff < -time.Second || diff > time.Second {
		t.Errorf("expected token TTL ~%v, got %v", wantTTL, gotTTL)
	}
}

func TestGenerateToken(t *testing.T) {
	secret := []byte("test-secret-key")
	h := &Handler{JWTSecret: secret}

	// NOTE: assumes maxEmailLength is well under ~500 chars (e.g. RFC 5321's
	// 254 cap). Adjust the repeat count if your constant is larger.
	tooLongEmail := strings.Repeat("a", 500) + "@example.com"

	tests := []struct {
		name          string
		req           *tokenpb.GenerateTokenRequest
		wantCode      codes.Code
		expectedError string
		validate      func(t *testing.T, resp *tokenpb.GenerateTokenResponse)
	}{
		{
			name: "succeeds with a valid admin id and email",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    42,
				AdminEmail: "admin@example.com",
			},
			validate: func(t *testing.T, resp *tokenpb.GenerateTokenResponse) {
				assertValidToken(t, resp, secret, 42, "admin@example.com")
			},
		},
		{
			name:          "fails when request is nil",
			req:           nil,
			wantCode:      codes.InvalidArgument,
			expectedError: "request cannot be nil",
		},
		{
			name: "fails when admin id is zero",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    0,
				AdminEmail: "admin@example.com",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_id is required and must be greater than 0",
		},
		{
			name: "fails when admin id is negative",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    -5,
				AdminEmail: "admin@example.com",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_id is required and must be greater than 0",
		},
		{
			name: "fails when admin email is empty",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_email is required",
		},
		{
			name: "fails when admin email is only whitespace",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "   ",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_email is required",
		},
		{
			name: "fails when admin email exceeds max length",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: tooLongEmail,
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_email exceeds maximum length",
		},
		{
			name: "fails when admin email contains internal whitespace",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "admin user@example.com",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_email cannot contain whitespace",
		},
		{
			name: "fails when admin email is malformed",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "not-an-email",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_email is invalid",
		},
		{
			name: "fails when admin email has no domain dot",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "admin@localhost",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_email is invalid",
		},
		{
			name: "fails and aggregates multiple validation errors",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    0,
				AdminEmail: "",
			},
			wantCode:      codes.InvalidArgument,
			expectedError: "admin_id is required and must be greater than 0; admin_email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := h.GenerateToken(context.Background(), tt.req)

			if tt.expectedError != "" {
				assertGRPCError(t, err, tt.wantCode, tt.expectedError)
				if resp != nil {
					t.Errorf("expected nil response on error, got %+v", resp)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tt.validate != nil {
				tt.validate(t, resp)
			}
		})
	}
}
