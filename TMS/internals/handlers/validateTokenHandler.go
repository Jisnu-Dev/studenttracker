package handlers

import (
	"context"
	"errors"
	"log/slog"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) ValidateToken(ctx context.Context, req *tokenpb.ValidateTokenRequest) (*tokenpb.ValidateTokenResponse, error) {
	if err := ValidateValidateTokenReq(req); err != nil {
		return nil, err
	}

	parseWithClaims := h.TokenParser
	if parseWithClaims == nil {
		parseWithClaims = jwt.ParseWithClaims
	}

	token, err := parseWithClaims(req.GetToken(), &Claims{}, func(token *jwt.Token) (interface{}, error) {
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
		if errors.Is(err, jwt.ErrTokenMalformed) ||
			errors.Is(err, jwt.ErrTokenSignatureInvalid) ||
			errors.Is(err, jwt.ErrTokenExpired) ||
			errors.Is(err, jwt.ErrTokenNotValidYet) ||
			errors.Is(err, jwt.ErrTokenInvalidClaims) {

			slog.Warn("invalid token",
				slog.String("function", "ValidateToken"),
				slog.Any("error", err),
			)
			return ValidateTokenResponse(false, 0, ""), nil
		}

		slog.Error("internal server error during token parsing",
			slog.String("function", "ValidateToken"),
			slog.Any("error", err),
		)
		return nil, ErrValidateTokenFailed
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return ValidateTokenResponse(true, claims.AdminID, claims.AdminEmail), nil
	}

	return ValidateTokenResponse(false, 0, ""), nil
}
