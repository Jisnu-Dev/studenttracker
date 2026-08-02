package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	AdminID    int64  `json:"adminID"`
	AdminEmail string `json:"adminEmail"`
	jwt.RegisteredClaims
}

type Service struct {
	JWTSecret []byte
}

func NewService(JWTSecret string) *Service {
	return &Service{
		JWTSecret: []byte(JWTSecret),
	}
}

func (s *Service) GenerateToken(adminID int64, adminEmail string) (string, error) {
	claims := Claims{
		AdminID:    adminID,
		AdminEmail: adminEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.JWTSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *Service) ValidateToken(tokenString string) (bool, int64, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.JWTSecret, nil
	})
	if err != nil {
		return false, 0, "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return true, claims.AdminID, claims.AdminEmail, nil
	}

	return false, 0, "", fmt.Errorf("invalid token")
}
