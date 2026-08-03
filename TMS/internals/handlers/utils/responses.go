package utils

import (
	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
)

// TokenSuccessResponse constructs a successful GenerateTokenResponse.
func TokenSuccessResponse(token string) *tokenpb.GenerateTokenResponse {
	return &tokenpb.GenerateTokenResponse{
		Token: token,
	}
}

// ValidTokenResponse constructs a valid ValidateTokenResponse.
func ValidTokenResponse(adminID int64, adminEmail string) *tokenpb.ValidateTokenResponse {
	return &tokenpb.ValidateTokenResponse{
		IsValid:    true,
		AdminId:    adminID,
		AdminEmail: adminEmail,
	}
}

// InvalidTokenResponse constructs an invalid ValidateTokenResponse.
func InvalidTokenResponse() *tokenpb.ValidateTokenResponse {
	return &tokenpb.ValidateTokenResponse{
		IsValid: false,
	}
}
