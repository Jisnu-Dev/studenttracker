package grpcHandler

import (
	"errors"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	AdminID    int64  `json:"adminID"`
	AdminEmail string `json:"adminEmail"`
	jwt.RegisteredClaims
}

var (
	// Fallback Function Execution Errors
	ErrGenerateTokenFailed = errors.New("failed to sign jwt token")
	ErrValidateTokenFailed = errors.New("unable to validate token")

	// Validation Errors
	ErrRequestNil             = errors.New("request cannot be nil")
	ErrAdminIDInvalid         = errors.New("admin_id is required and must be greater than 0")
	ErrAdminEmailRequired     = errors.New("admin_email is required")
	ErrAdminEmailTooLong      = errors.New("admin_email exceeds maximum length")
	ErrAdminEmailNoWhitespace = errors.New("admin_email cannot contain whitespace")
	ErrAdminEmailInvalid      = errors.New("admin_email is invalid")
	ErrTokenRequired          = errors.New("token is required")
	ErrTokenNoWhitespace      = errors.New("token cannot contain whitespace")
	ErrTokenMalformed         = errors.New("malformed token structure")
)

type Handler struct {
	JWTSecret   any
	tokenpb.UnimplementedTokenServiceServer
}

func NewHandler(JWTSecret string) *Handler {
	return &Handler{
		JWTSecret:   []byte(JWTSecret),
	}
}
