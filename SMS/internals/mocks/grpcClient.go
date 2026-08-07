package mocks

import (
	"context"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockToken is a standard mock JWT token string used across tests.
const MockToken = "signed.jwt.token"

// MockTokenServiceClient is a mock implementation of the TMS gRPC token client.
type MockTokenServiceClient struct {
	// Error simulation triggers
	GenerateErr GrpcOpError
	ValidateErr GrpcOpError

	// Mock return data configuration
	Token      string
	AdminID    int64
	AdminEmail string
}

// GenerateToken mocks a token generation request to the TMS service.
func (m *MockTokenServiceClient) GenerateToken(ctx context.Context, req *tokenpb.GenerateTokenRequest, _ ...grpc.CallOption) (*tokenpb.GenerateTokenResponse, error) {
	switch m.GenerateErr {
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, "failed to generate token")
	case GrpcOpUnavailable:
		return nil, status.Error(codes.Unavailable, "TMS service unavailable")
	}

	tok := m.Token
	if tok == "" {
		tok = MockToken
	}

	return &tokenpb.GenerateTokenResponse{Token: tok}, nil
}

// ValidateToken mocks a token validation request to the TMS service.
func (m *MockTokenServiceClient) ValidateToken(ctx context.Context, req *tokenpb.ValidateTokenRequest, _ ...grpc.CallOption) (*tokenpb.ValidateTokenResponse, error) {
	switch m.ValidateErr {
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, "token validation service unavailable")
	case GrpcOpInvalidToken:
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
	case GrpcOpUnavailable:
		return nil, status.Error(codes.Unavailable, "TMS service unavailable")
	}

	return &tokenpb.ValidateTokenResponse{
		IsValid:    true,
		AdminID:    m.AdminID,
		AdminEmail: m.AdminEmail,
	}, nil
}
