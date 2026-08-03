package utils

import (
	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
)

// GenerateTokenResponse constructs a GenerateTokenResponse proto message.
func GenerateTokenResponse(token string) *tokenpb.GenerateTokenResponse {
	return &tokenpb.GenerateTokenResponse{
		Token: token,
	}
}

// ValidateTokenResponse constructs a ValidateTokenResponse proto message.
func ValidateTokenResponse(isValid bool, adminID int64, adminEmail string) *tokenpb.ValidateTokenResponse {
	return &tokenpb.ValidateTokenResponse{
		IsValid:    isValid,
		AdminId:    adminID,
		AdminEmail: adminEmail,
	}
}
