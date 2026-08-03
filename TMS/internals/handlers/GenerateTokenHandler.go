package handlers

import (
	"context"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/Jisnu-Dev/TMS/internals/handlers/utils"
	"google.golang.org/grpc/codes"
)

func (h *Handler) GenerateToken(ctx context.Context, req *tokenpb.GenerateTokenRequest) (*tokenpb.GenerateTokenResponse, error) {
	if err := utils.ValidateGenerateTokenReq(req); err != nil {
		return nil, err
	}

	token, err := h.service.GenerateToken(req.GetAdminID(), req.GetAdminEmail())
	if err != nil {
		return nil, utils.GrpcError(codes.Internal, err.Error())
	}

	return utils.GenerateTokenResponse(token), nil
}
