package clients

import (
	"context"
	"fmt"
	"time"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

func (c *TokenClient) GenerateToken(ctx context.Context, adminID int64, adminEmail string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.GenerateToken(ctx, &tokenpb.GenerateTokenRequest{
		AdminID:    adminID,
		AdminEmail: adminEmail,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return resp.GetToken(), nil
}

func (c *TokenClient) ValidateToken(ctx context.Context, token string) (bool, int64, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.ValidateToken(ctx, &tokenpb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return false, 0, "", fmt.Errorf("failed to validate token: %w", err)
	}

	return resp.GetIsValid(), resp.GetAdminID(), resp.GetAdminEmail(), nil
}
