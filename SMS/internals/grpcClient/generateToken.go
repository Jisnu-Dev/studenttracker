package grpcClient

import (
	"context"
	"time"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
)

func (c *TokenClient) GenerateToken(ctx context.Context, adminID int64, adminEmail string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.GenerateToken(ctx, &tokenpb.GenerateTokenRequest{
		AdminID:    adminID,
		AdminEmail: adminEmail,
	})
	if err != nil {
		return "", err
	}

	return resp.GetToken(), nil
}
