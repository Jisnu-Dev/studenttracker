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

type Handler struct {
	JWTSecret []byte
	tokenpb.UnimplementedTokenServiceServer
}

func NewHandler(JWTSecret string) *Handler {
	return &Handler{
		JWTSecret: []byte(JWTSecret),
	}
}
