package mocks

import (
	"context"
	"errors"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
	"google.golang.org/grpc"
)

type MockTokenServiceClient struct {
	// 1. Error simulation triggers
	GenerateTokenError MockOpError
	ValidateTokenError MockOpError

	// 2. Mock return data configuration
	Token      string
	AdminID    int64
	AdminEmail string

	// 3. Captured arguments (for verification in tests)
	CapturedAdminID    int64
	CapturedAdminEmail string
	CapturedToken      string
}

func (m *MockTokenServiceClient) GenerateToken(ctx context.Context, req *tokenpb.GenerateTokenRequest, _ ...grpc.CallOption) (*tokenpb.GenerateTokenResponse, error) {
	m.CapturedAdminID = req.GetAdminID()
	m.CapturedAdminEmail = req.GetAdminEmail()

	switch m.GenerateTokenError {
	case OpInternalError:
		return nil, errors.New("rpc error: failed to generate token")
	}

	if m.Token != "" {
		return &tokenpb.GenerateTokenResponse{Token: m.Token}, nil
	}
	return &tokenpb.GenerateTokenResponse{Token: "mock.jwt.token"}, nil
}

func (m *MockTokenServiceClient) ValidateToken(ctx context.Context, req *tokenpb.ValidateTokenRequest, _ ...grpc.CallOption) (*tokenpb.ValidateTokenResponse, error) {
	m.CapturedToken = req.GetToken()

	switch m.ValidateTokenError {
	case OpInternalError:
		return nil, errors.New("rpc error: token validation service unavailable")
	case OpInvalidToken, OpNotFound:
		return &tokenpb.ValidateTokenResponse{
			IsValid: false,
		}, nil
	}

	adminID := int64(1)
	if m.AdminID != 0 {
		adminID = m.AdminID
	}

	email := "admin@example.com"
	if m.AdminEmail != "" {
		email = m.AdminEmail
	}

	return &tokenpb.ValidateTokenResponse{
		IsValid:    true,
		AdminID:    adminID,
		AdminEmail: email,
	}, nil
}
