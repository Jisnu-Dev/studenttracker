package grpcHandler

import (
	"context"
	"log/slog"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) ValidateToken(ctx context.Context, req *tokenpb.ValidateTokenRequest) (*tokenpb.ValidateTokenResponse, error) {
	if err := ValidateTokenReq(req); err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(req.GetToken(), &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Warn("unexpected signing method",
				slog.String("function", "ValidateToken"),
				slog.Any("alg", token.Header["alg"]),
			)
			return nil, jwt.ErrSignatureInvalid
		}
		return h.JWTSecret, nil
	})

	if err != nil {
		slog.Warn("invalid token",
			slog.String("function", "ValidateToken"),
			slog.Any("error", err),
		)
		return ValidateTokenResponse(false, 0, ""), nil
	}

	claims := token.Claims.(*Claims)
	return ValidateTokenResponse(true, claims.AdminID, claims.AdminEmail), nil
}
