package utils

import (
	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
)

func GenerateTokenResponse(token string) *tokenpb.GenerateTokenResponse {
	return &tokenpb.GenerateTokenResponse{
		Token: token,
	}
}

func ValidateTokenResponse(isValid bool, adminID int64, adminEmail string) *tokenpb.ValidateTokenResponse {
	return &tokenpb.ValidateTokenResponse{
		IsValid:    isValid,
		AdminId:    adminID,
		AdminEmail: adminEmail,
	}
}
