package grpcClient

import (
	"context"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ValidateTokenResponse struct {
	IsValid    bool
	AdminID    int64
	AdminEmail string
}

type TokenClientInterface interface {
	GenerateToken(ctx context.Context, adminID int64, adminEmail string) (string, error)
	ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error)
}

type TokenClient struct {
	client tokenpb.TokenServiceClient
	conn   *grpc.ClientConn
}

func NewTokenClient(TMSAddress string) (*TokenClient, error) {
	conn, err := grpc.NewClient(TMSAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &TokenClient{
		client: tokenpb.NewTokenServiceClient(conn),
		conn:   conn,
	}, nil
}

// NewTokenClientForTest creates a TokenClient with an injected gRPC client,
// used only in tests to avoid needing a real network connection.
func NewTokenClientForTest(c tokenpb.TokenServiceClient) *TokenClient {
	return &TokenClient{client: c}
}
