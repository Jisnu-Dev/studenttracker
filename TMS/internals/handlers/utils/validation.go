package utils

import (
	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
)

// ValidateGenerateTokenReq validates the incoming GenerateToken request.
func ValidateGenerateTokenReq(req *tokenpb.GenerateTokenRequest) error {
	if req == nil {
		return ErrInvalidArgument("request cannot be nil")
	}
	if req.GetAdminID() == 0 {
		return ErrInvalidArgument("admin_id is required and must be greater than 0")
	}
	if req.GetAdminEmail() == "" {
		return ErrInvalidArgument("admin_email is required")
	}
	return nil
}

// ValidateValidateTokenReq validates the incoming ValidateToken request.
func ValidateValidateTokenReq(req *tokenpb.ValidateTokenRequest) error {
	if req == nil || req.GetToken() == "" {
		return ErrInvalidArgument("token is required")
	}
	return nil
}
