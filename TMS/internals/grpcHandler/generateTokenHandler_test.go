package grpcHandler

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
func assertGRPCError(t *testing.T, err error, expectedCode codes.Code, expectedError string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected a gRPC status error, got %v", err)
	}

	if st.Code() != expectedCode {
		t.Errorf("expected code %v, got %v (message: %q)", expectedCode, st.Code(), st.Message())
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

	tooLongEmail := strings.Repeat("a", 500) + "@example.com"

	tests := []struct {
		name          string
		handler       *Handler
		req           *tokenpb.GenerateTokenRequest
		expectedCode  codes.Code
		expectedError string
		validate      func(t *testing.T, resp *tokenpb.GenerateTokenResponse)
	}{
		{
			name: "generate token successful with valid admin id and email",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    42,
				AdminEmail: "admin@example.com",
			},
			expectedCode: codes.OK,
			validate: func(t *testing.T, resp *tokenpb.GenerateTokenResponse) {
				assertValidToken(t, resp, secret, 42, "admin@example.com")
			},
		},
		{
			name:          "generate token fails when request is nil",
			req:           nil,
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrRequestNil.Error(),
		},
		{
			name: "generate token fails when admin id is zero",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    0,
				AdminEmail: "admin@example.com",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "admin_id is required and must be greater than 0",
		},
		{
			name: "generate token fails when admin id is negative",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    -5,
				AdminEmail: "admin@example.com",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrAdminIDInvalid.Error(),
		},
		{
			name: "generate token fails when admin email is empty",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "admin_email is required",
		},
		{
			name: "generate token fails when admin email is only whitespace",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "   ",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrAdminEmailRequired.Error(),
		},
		{
			name: "generate token fails when admin email exceeds max length",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: tooLongEmail,
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrAdminEmailTooLong.Error(),
		},
		{
			name: "generate token fails when admin email contains internal whitespace",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "admin user@example.com",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrAdminEmailNoWhitespace.Error(),
		},
		{
			name: "generate token fails when admin email is malformed",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "not-an-email",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "admin_email is invalid",
		},
		{
			name: "generate token fails when admin email has no domain dot",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "admin@localhost",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrAdminEmailInvalid.Error(),
		},
		{
			name: "generate token fails and aggregates multiple validation errors",
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    0,
				AdminEmail: "",
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrAdminIDInvalid.Error() + "; " + ErrAdminEmailRequired.Error(),
		},
		{
			name:    "generate token fails when jwt token signing fails",
			handler: &Handler{JWTSecret: "invalid-key-type"},
			req: &tokenpb.GenerateTokenRequest{
				AdminID:    42,
				AdminEmail: "admin@example.com",
			},
			expectedCode:  codes.Internal,
			expectedError: ErrGenerateTokenFailed.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetHandler := h
			if tt.handler != nil {
				targetHandler = tt.handler
			}

			resp, err := targetHandler.GenerateToken(context.Background(), tt.req)

			if tt.expectedCode != codes.OK {
				assertGRPCError(t, err, tt.expectedCode, tt.expectedError)
				if resp != nil {
					t.Errorf("expected nil response on error, got %+v", resp)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.validate != nil {
				tt.validate(t, resp)
			}
		})
	}
}
