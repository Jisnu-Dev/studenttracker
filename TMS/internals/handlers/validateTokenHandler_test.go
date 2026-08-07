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
		name          string
		handler       *Handler
		req           *tokenpb.ValidateTokenRequest
		token         string
		expectedCode  codes.Code
		expectedError string

		expectedIsValid bool
		expectedAdminID int64
		expectedEmail   string
	}{
		{
			name:          "validate token fails when request is nil",
			req:           nil,
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrRequestNil.Error(),
		},
		{
			name:          "validate token fails when token is empty",
			token:         "",
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrTokenRequired.Error(),
		},
		{
			name:          "validate token fails when token is only whitespace",
			token:         "   ",
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrTokenRequired.Error(),
		},
		{
			name:          "validate token fails when token contains internal whitespace",
			token:         "abc def.ghi.jkl",
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrTokenNoWhitespace.Error(),
		},
		{
			name:          "validate token fails when token structure has wrong number of segments",
			token:         "abcdef",
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrTokenMalformed.Error(),
		},
		{
			name:          "validate token fails and aggregates multiple structural errors",
			token:         "abc def",
			expectedCode:  codes.InvalidArgument,
			expectedError: ErrTokenNoWhitespace.Error() + "; " + ErrTokenMalformed.Error(),
		},
		{
			name:            "validate token returns invalid for structurally-valid but garbage token",
			token:           structurallyValidButGarbageToken,
			expectedCode:    codes.OK,
			expectedIsValid: false,
		},
		{
			name:            "validate token returns invalid for token signed with wrong secret",
			token:           wrongSigToken,
			expectedCode:    codes.OK,
			expectedIsValid: false,
		},
		{
			name:            "validate token returns invalid for expired token",
			token:           expiredToken,
			expectedCode:    codes.OK,
			expectedIsValid: false,
		},
		{
			name:            "validate token returns invalid for not-yet-valid token",
			token:           notYetValidToken,
			expectedCode:    codes.OK,
			expectedIsValid: false,
		},
		{
			name:            "validate token returns invalid for algorithm mismatch token",
			token:           algMismatchToken,
			expectedCode:    codes.OK,
			expectedIsValid: false,
		},
		{
			name: "validate token returns internal error when parser fails unexpectedly",
			handler: &Handler{
				JWTSecret: secret,
				TokenParser: func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, _ ...jwt.ParserOption) (*jwt.Token, error) {
					return nil, errors.New("unexpected parser failure")
				},
			},
			token:         "valid.dummy.token",
			expectedCode:  codes.Internal,
			expectedError: ErrValidateTokenFailed.Error(),
		},
		{
			name: "validate token returns invalid when parsed token is marked invalid with nil error",
			handler: &Handler{
				JWTSecret: secret,
				TokenParser: func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, _ ...jwt.ParserOption) (*jwt.Token, error) {
					return &jwt.Token{Valid: false, Claims: &Claims{AdminID: 7, AdminEmail: "admin@example.com"}}, nil
				},
			},
			token:           "valid.dummy.token",
			expectedCode:    codes.OK,
			expectedIsValid: false,
		},
		{
			name: "validate token returns invalid when parsed token claims is not *Claims",
			handler: &Handler{
				JWTSecret: secret,
				TokenParser: func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, _ ...jwt.ParserOption) (*jwt.Token, error) {
					return &jwt.Token{Valid: true, Claims: jwt.MapClaims{"adminID": 7}}, nil
				},
			},
			token:           "valid.dummy.token",
			expectedCode:    codes.OK,
			expectedIsValid: false,
		},
		{
			name:            "validate token successful for valid token and returns its claims",
			token:           validToken,
			expectedCode:    codes.OK,
			expectedIsValid: true,
			expectedAdminID: 7,
			expectedEmail:   "admin@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.req
			if req == nil && tt.name != "validate token fails when request is nil" {
				req = &tokenpb.ValidateTokenRequest{Token: tt.token}
			}

			targetHandler := h
			if tt.handler != nil {
				targetHandler = tt.handler
			}

			resp, err := targetHandler.ValidateToken(context.Background(), req)

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
			if resp == nil {
				t.Fatal("expected a non-nil response")
			}
			if resp.GetIsValid() != tt.expectedIsValid {
				t.Errorf("expected IsValid %v, got %v", tt.expectedIsValid, resp.GetIsValid())
			}
			if resp.GetAdminId() != tt.expectedAdminID {
				t.Errorf("expected AdminId %d, got %d", tt.expectedAdminID, resp.GetAdminId())
			}
			if resp.GetAdminEmail() != tt.expectedEmail {
				t.Errorf("expected AdminEmail %q, got %q", tt.expectedEmail, resp.GetAdminEmail())
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
