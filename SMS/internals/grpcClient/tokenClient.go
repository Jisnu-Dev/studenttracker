package grpcClient

import (
	"context"
	"fmt"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TokenClientInterface interface {
	GenerateToken(ctx context.Context, adminID int64, adminEmail string) (string, error)
	ValidateToken(ctx context.Context, token string) (bool, int64, string, error)
}

type TokenClient struct {
	client tokenpb.TokenServiceClient
	conn   *grpc.ClientConn
}

func NewTokenClient(TMSAddress string) (*TokenClient, error) {
	conn, err := grpc.NewClient(TMSAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to TMS: %w", err)
	}

	return &TokenClient{
		client: tokenpb.NewTokenServiceClient(conn),
		conn:   conn,
	}, nil
}

// newTokenClientForTest creates a TokenClient with an injected gRPC client,
// used only in tests to avoid needing a real network connection.
func newTokenClientForTest(c tokenpb.TokenServiceClient) *TokenClient {
	return &TokenClient{client: c}
}
