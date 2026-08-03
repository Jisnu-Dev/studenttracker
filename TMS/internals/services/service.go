package services

import (
	"log/slog"
	"time"

	serviceErrors "github.com/Jisnu-Dev/TMS/internals/services/errors"
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
		slog.Error("failed to sign jwt token",
			slog.String("function", "GenerateToken"),
			slog.Int64("adminID", adminID),
			slog.String("adminEmail", adminEmail),
			slog.Any("error", err),
		)
		return "", serviceErrors.ErrGenerateTokenFailed
	}

	return tokenString, nil
}

func (s *Service) ValidateToken(tokenString string) (bool, int64, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Warn("unexpected signing method",
				slog.String("function", "ValidateToken"),
				slog.Any("alg", token.Header["alg"]),
			)
			return nil, serviceErrors.ErrSigningMethodMismatch
		}
		return s.JWTSecret, nil
	})
	if err != nil {
		slog.Warn("token parsing failed",
			slog.String("function", "ValidateToken"),
			slog.Any("error", err),
		)
		return false, 0, "", serviceErrors.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return true, claims.AdminID, claims.AdminEmail, nil
	}

	return false, 0, "", serviceErrors.ErrInvalidToken
}
