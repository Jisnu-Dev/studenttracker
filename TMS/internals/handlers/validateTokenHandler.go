package handlers

import (
	"context"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/Jisnu-Dev/TMS/internals/handlers/utils"
)

func (h *Handler) ValidateToken(ctx context.Context, req *tokenpb.ValidateTokenRequest) (*tokenpb.ValidateTokenResponse, error) {
	if err := utils.ValidateValidateTokenReq(req); err != nil {
		return nil, err
	}

	isValid, adminId, adminEmail, err := h.service.ValidateToken(req.GetToken())
	if err != nil || !isValid {
		return utils.InvalidTokenResponse(), nil
	}

	return utils.ValidTokenResponse(adminId, adminEmail), nil
}
