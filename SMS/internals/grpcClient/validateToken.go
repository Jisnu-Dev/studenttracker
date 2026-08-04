package grpcClient

import (
	"context"
	"fmt"
	"time"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
)

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
