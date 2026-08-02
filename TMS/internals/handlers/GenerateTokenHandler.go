package handlers

import (
	"context"
	"fmt"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GenerateToken(ctx context.Context, req *tokenpb.GenerateTokenRequest) (*tokenpb.GenerateTokenResponse, error) {
	if req.GetAdminID() == 0 || req.GetAdminEmail() == "" {
		return nil, fmt.Errorf("admin id and email are required")
	}

	token, err := h.service.GenerateToken(req.GetAdminID(), req.GetAdminEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}
	return &tokenpb.GenerateTokenResponse{
		Token: token,
	}, nil
}

func (h *Handler) ValidateToken(ctx context.Context, req *tokenpb.ValidateTokenRequest) (*tokenpb.ValidateTokenResponse, error) {
	if req.GetToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	isValid, adminId, adminEmail, err := h.service.ValidateToken(req.GetToken())
	if err != nil || !isValid {
		return &tokenpb.ValidateTokenResponse{
			IsValid: false,
		}, nil
	}

	return &tokenpb.ValidateTokenResponse{
		IsValid:    isValid,
		AdminId:    adminId,
		AdminEmail: adminEmail,
	}, nil
}
