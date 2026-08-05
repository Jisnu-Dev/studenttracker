package handlers

import (
	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	AdminID    int64  `json:"adminID"`
	AdminEmail string `json:"adminEmail"`
	jwt.RegisteredClaims
}

type TokenParser func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error)

type Handler struct {
	JWTSecret   any
	TokenParser TokenParser
	tokenpb.UnimplementedTokenServiceServer
}

func NewHandler(JWTSecret string) *Handler {
	return &Handler{
		JWTSecret:   []byte(JWTSecret),
		TokenParser: jwt.ParseWithClaims,
	}
}
