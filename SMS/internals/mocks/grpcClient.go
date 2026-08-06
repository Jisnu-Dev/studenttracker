package mocks

import (
	"context"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockTokenServiceClient is a mock implementation of the TMS gRPC token client.
type MockTokenServiceClient struct {
	// 1. Error simulation triggers
	GenerateErr GrpcOpError
	ValidateErr GrpcOpError

	// 2. Mock return data configuration
	Token      string
	AdminID    int64
	AdminEmail string

	// 3. Captured arguments (for verification in tests)
	CapturedAdminID    int64
	CapturedAdminEmail string
	CapturedToken      string
}

// GenerateToken mocks a token generation request to the TMS service.
func (m *MockTokenServiceClient) GenerateToken(ctx context.Context, req *tokenpb.GenerateTokenRequest, _ ...grpc.CallOption) (*tokenpb.GenerateTokenResponse, error) {
	m.CapturedAdminID = req.GetAdminID()
	m.CapturedAdminEmail = req.GetAdminEmail()

	switch m.GenerateErr {
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	if m.Token != "" {
		return &tokenpb.GenerateTokenResponse{Token: m.Token}, nil
	}
	return &tokenpb.GenerateTokenResponse{Token: "mock.jwt.token"}, nil
}

// ValidateToken mocks a token validation request to the TMS service.
func (m *MockTokenServiceClient) ValidateToken(ctx context.Context, req *tokenpb.ValidateTokenRequest, _ ...grpc.CallOption) (*tokenpb.ValidateTokenResponse, error) {
	m.CapturedToken = req.GetToken()

	switch m.ValidateErr {
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, "token validation service unavailable")
	case GrpcOpInvalidToken:
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
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
