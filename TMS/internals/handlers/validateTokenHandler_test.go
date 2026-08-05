package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
)

func mustSignToken(t *testing.T, secret []byte, claims Claims) string {
	t.Helper()

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString(secret)
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}
	return s
}

func mustSignNoneAlgToken(t *testing.T, claims Claims) string {
	t.Helper()

	tok := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	s, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("failed to sign none-alg test token: %v", err)
	}
	return s
}

func TestValidateToken(t *testing.T) {
	secret := []byte("test-secret-key")
	h := &Handler{JWTSecret: secret}

	now := time.Now()

	validClaims := Claims{
		AdminID:    7,
		AdminEmail: "admin@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Hour)),
		},
	}

	expiredClaims := Claims{
		AdminID:    7,
		AdminEmail: "admin@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now.Add(-10 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(now.Add(-5 * time.Hour)),
		},
	}

	notYetValidClaims := Claims{
		AdminID:    7,
		AdminEmail: "admin@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(1 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Hour)),
		},
	}

	validToken := mustSignToken(t, secret, validClaims)
	expiredToken := mustSignToken(t, secret, expiredClaims)
	notYetValidToken := mustSignToken(t, secret, notYetValidClaims)
	wrongSigToken := mustSignToken(t, []byte("a-different-secret"), validClaims)
	algMismatchToken := mustSignNoneAlgToken(t, validClaims)

	const structurallyValidButGarbageToken = "abc.!!!.def"

	tests := []struct {
		name    string
		handler *Handler
		req     *tokenpb.ValidateTokenRequest
		token   string // used instead of req when req is nil, for readability

		wantGRPCErr   bool
		wantCode      codes.Code
		expectedError string

		wantInternalErr bool // expects the bare ErrValidateTokenFailed sentinel

		wantIsValid    bool
		wantAdminID    int64
		wantAdminEmail string
	}{
		{
			name:          "fails when request is nil",
			req:           nil,
			wantGRPCErr:   true,
			wantCode:      codes.InvalidArgument,
			expectedError: "request cannot be nil",
		},
		{
			name:          "fails when token is empty",
			token:         "",
			wantGRPCErr:   true,
			wantCode:      codes.InvalidArgument,
			expectedError: "token is required",
		},
		{
			name:          "fails when token is only whitespace",
			token:         "   ",
			wantGRPCErr:   true,
			wantCode:      codes.InvalidArgument,
			expectedError: "token is required",
		},
		{
			name:          "fails when token contains internal whitespace",
			token:         "abc def.ghi.jkl", // 2 dots, so only the whitespace check should trip
			wantGRPCErr:   true,
			wantCode:      codes.InvalidArgument,
			expectedError: "token cannot contain whitespace",
		},
		{
			name:          "fails when token structure has the wrong number of segments",
			token:         "abcdef", // 0 dots
			wantGRPCErr:   true,
			wantCode:      codes.InvalidArgument,
			expectedError: "malformed token structure",
		},
		{
			name:          "fails and aggregates multiple structural errors",
			token:         "abc def", // 0 dots AND contains whitespace
			wantGRPCErr:   true,
			wantCode:      codes.InvalidArgument,
			expectedError: "token cannot contain whitespace; malformed token structure",
		},
		{
			name:           "returns invalid for a structurally-valid but garbage token",
			token:          structurallyValidButGarbageToken,
			wantIsValid:    false,
			wantAdminID:    0,
			wantAdminEmail: "",
		},
		{
			name:           "returns invalid for a token signed with the wrong secret",
			token:          wrongSigToken,
			wantIsValid:    false,
			wantAdminID:    0,
			wantAdminEmail: "",
		},
		{
			name:           "returns invalid for an expired token",
			token:          expiredToken,
			wantIsValid:    false,
			wantAdminID:    0,
			wantAdminEmail: "",
		},
		{
			name:           "returns invalid for a not-yet-valid token",
			token:          notYetValidToken,
			wantIsValid:    false,
			wantAdminID:    0,
			wantAdminEmail: "",
		},
		{
			name:            "BUG: algorithm mismatch is misclassified as an internal error",
			token:           algMismatchToken,
			wantInternalErr: true,
		},
		{
			name: "returns invalid when parsed token is marked invalid with nil error",
			handler: &Handler{
				JWTSecret: secret,
				TokenParser: func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, _ ...jwt.ParserOption) (*jwt.Token, error) {
					return &jwt.Token{Valid: false, Claims: &Claims{AdminID: 7, AdminEmail: "admin@example.com"}}, nil
				},
			},
			token:          "valid.dummy.token",
			wantIsValid:    false,
			wantAdminID:    0,
			wantAdminEmail: "",
		},
		{
			name: "returns invalid when parsed token claims is not *Claims",
			handler: &Handler{
				JWTSecret: secret,
				TokenParser: func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, _ ...jwt.ParserOption) (*jwt.Token, error) {
					return &jwt.Token{Valid: true, Claims: jwt.MapClaims{"adminID": 7}}, nil
				},
			},
			token:          "valid.dummy.token",
			wantIsValid:    false,
			wantAdminID:    0,
			wantAdminEmail: "",
		},
		{
			name:           "succeeds for a valid token and returns its claims",
			token:          validToken,
			wantIsValid:    true,
			wantAdminID:    7,
			wantAdminEmail: "admin@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.req
			if req == nil && tt.name != "fails when request is nil" {
				req = &tokenpb.ValidateTokenRequest{Token: tt.token}
			}

			targetHandler := h
			if tt.handler != nil {
				targetHandler = tt.handler
			}

			resp, err := targetHandler.ValidateToken(context.Background(), req)

			switch {
			case tt.wantGRPCErr:
				assertGRPCError(t, err, tt.wantCode, tt.expectedError)
				if resp != nil {
					t.Errorf("expected nil response on error, got %+v", resp)
				}

			case tt.wantInternalErr:
				if !errors.Is(err, ErrValidateTokenFailed) {
					t.Errorf("expected ErrValidateTokenFailed, got %v", err)
				}
				if resp != nil {
					t.Errorf("expected nil response on internal error, got %+v", resp)
				}

			default:
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if resp == nil {
					t.Fatal("expected a non-nil response")
				}
				if resp.GetIsValid() != tt.wantIsValid {
					t.Errorf("expected IsValid %v, got %v", tt.wantIsValid, resp.GetIsValid())
				}
				if resp.GetAdminId() != tt.wantAdminID {
					t.Errorf("expected AdminId %d, got %d", tt.wantAdminID, resp.GetAdminId())
				}
				if resp.GetAdminEmail() != tt.wantAdminEmail {
					t.Errorf("expected AdminEmail %q, got %q", tt.wantAdminEmail, resp.GetAdminEmail())
				}
			}
		})
	}
}

func TestNewHandler(t *testing.T) {
	h := NewHandler("my-secret-key")
	if string(h.JWTSecret.([]byte)) != "my-secret-key" {
		t.Errorf("expected secret %q, got %v", "my-secret-key", h.JWTSecret)
	}
	if h.TokenParser == nil {
		t.Error("expected TokenParser to be initialized")
	}
}
