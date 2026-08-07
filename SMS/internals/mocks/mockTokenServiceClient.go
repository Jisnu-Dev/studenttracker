package mocks

import (
	"context"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
	"google.golang.org/grpc"
)

const MockToken = "signed.jwt.token"

type MockTokenServiceClient struct {
	GenerateTokenError GrpcOpError
	ValidateTokenError GrpcOpError

	ReturnInvalidToken bool
	AdminID            int64
	AdminEmail         string

	// Captured requests, for asserting on what was actually sent.
	LastGenerateTokenReq *tokenpb.GenerateTokenRequest
	LastValidateTokenReq *tokenpb.ValidateTokenRequest
}

func (m *MockTokenServiceClient) GenerateToken(ctx context.Context, req *tokenpb.GenerateTokenRequest, _ ...grpc.CallOption) (*tokenpb.GenerateTokenResponse, error) {
	m.LastGenerateTokenReq = req

	switch m.GenerateTokenError {
	case GrpcOpInternalError:
		return nil, ErrGenerateTokenInternal
	case GrpcOpUnavailable:
		return nil, ErrTMSUnavailable
	}

	return &tokenpb.GenerateTokenResponse{Token: MockToken}, nil
}

func (m *MockTokenServiceClient) ValidateToken(ctx context.Context, req *tokenpb.ValidateTokenRequest, _ ...grpc.CallOption) (*tokenpb.ValidateTokenResponse, error) {
	m.LastValidateTokenReq = req

	switch m.ValidateTokenError {
	case GrpcOpInternalError:
		return nil, ErrValidateTokenInternal
	case GrpcOpInvalidToken:
		return nil, ErrInvalidToken
	case GrpcOpUnavailable:
		return nil, ErrTMSUnavailable
	}

	if m.ReturnInvalidToken {
		return &tokenpb.ValidateTokenResponse{IsValid: false}, nil
	}

	return &tokenpb.ValidateTokenResponse{
		IsValid:    true,
		AdminID:    m.AdminID,
		AdminEmail: m.AdminEmail,
	}, nil
}
