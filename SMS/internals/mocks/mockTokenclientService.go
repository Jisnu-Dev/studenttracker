package mocks

import (
	"context"
	"errors"
)

type MockTokenClient struct {
	GenerateTokenError MockOpError
	ValidateTokenError MockOpError

	Token   string
	IsValid bool
	AdminID int64
	Email   string

	CapturedAdminID    int64
	CapturedAdminEmail string
	CapturedToken      string
}

func (m *MockTokenClient) GenerateToken(ctx context.Context, adminID int64, adminEmail string) (string, error) {
	m.CapturedAdminID = adminID
	m.CapturedAdminEmail = adminEmail

	switch m.GenerateTokenError {
	case OpInternalError:
		return "", errors.New("failed to generate token")
	}

	if m.Token != "" {
		return m.Token, nil
	}
	return "mock.jwt.token", nil
}

func (m *MockTokenClient) ValidateToken(ctx context.Context, token string) (bool, int64, string, error) {
	m.CapturedToken = token

	switch m.ValidateTokenError {
	case OpInternalError:
		return false, 0, "", errors.New("token validation service unavailable")
	case OpNotFound:
		return false, 0, "", nil
	}

	adminID := int64(1)
	if m.AdminID != 0 {
		adminID = m.AdminID
	}

	email := "admin@example.com"
	if m.Email != "" {
		email = m.Email
	}

	return true, adminID, email, nil
}
