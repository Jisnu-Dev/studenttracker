package grpcClient

import (
	"context"
	"time"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
)

func (c *TokenClient) ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.ValidateToken(ctx, &tokenpb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return &ValidateTokenResponse{
		IsValid:    resp.GetIsValid(),
		AdminID:    resp.GetAdminID(),
		AdminEmail: resp.GetAdminEmail(),
	}, nil
}
