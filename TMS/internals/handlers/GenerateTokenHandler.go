package handlers

import (
	"context"
	"log/slog"
	"time"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/Jisnu-Dev/TMS/internals/handlers/utils"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
)

func (h *Handler) GenerateToken(ctx context.Context, req *tokenpb.GenerateTokenRequest) (*tokenpb.GenerateTokenResponse, error) {
	if err := utils.ValidateGenerateTokenReq(req); err != nil {
		return nil, err
	}

	claims := Claims{
		AdminID:    req.GetAdminID(),
		AdminEmail: req.GetAdminEmail(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.JWTSecret)
	if err != nil {
		slog.Error("failed to sign jwt token",
			slog.String("function", "GenerateToken"),
			slog.Int64("adminID", req.GetAdminID()),
			slog.String("adminEmail", req.GetAdminEmail()),
			slog.Any("error", err),
		)
		return nil, utils.GrpcError(codes.Internal, ErrGenerateTokenFailed.Error())
	}

	return utils.GenerateTokenResponse(tokenString), nil
}
