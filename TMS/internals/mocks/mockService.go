package mocks

import services "github.com/Jisnu-Dev/TMS/internals/services/errors"

type MockService struct {
	// 1. Error simulation triggers
	GenerateTokenError MockOpError
	ValidateTokenError MockOpError

	// 2. Mock return data configuration
	Token      string
	IsValid    bool
	AdminID    int64
	AdminEmail string

	// 3. Captured arguments (for verification in tests)
	CapturedAdminID    int64
	CapturedAdminEmail string
	CapturedToken      string
}

func (m *MockService) GenerateToken(adminID int64, adminEmail string) (string, error) {
	m.CapturedAdminID = adminID
	m.CapturedAdminEmail = adminEmail

	switch m.GenerateTokenError {
	case OpInternalError:
		return "", services.ErrGenerateTokenFailed
	}

	if m.Token != "" {
		return m.Token, nil
	}
	return "mocked.jwt.token", nil
}

func (m *MockService) ValidateToken(tokenString string) (bool, int64, string, error) {
	m.CapturedToken = tokenString

	switch m.ValidateTokenError {
	case OpInvalidToken:
		return false, 0, "", services.ErrInvalidToken
	case OpSigningMethodMismatch:
		return false, 0, "", services.ErrSigningMethodMismatch
	case OpInternalError:
		return false, 0, "", services.ErrValidateTokenFailed
	}

	adminID := int64(1)
	if m.AdminID != 0 {
		adminID = m.AdminID
	}

	adminEmail := "admin@example.com"
	if m.AdminEmail != "" {
		adminEmail = m.AdminEmail
	}

	return true, adminID, adminEmail, nil
}
