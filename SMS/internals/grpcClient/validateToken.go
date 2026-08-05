package grpcClient

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
)

func (c *TokenClient) ValidateToken(ctx context.Context, token string) (bool, int64, string, error) {
	var errs []string
	trimmed := strings.TrimSpace(token)
	if trimmed == "" {
		errs = append(errs, "token is required")
	} else {
		if strings.ContainsAny(trimmed, " \t\n\r") {
			errs = append(errs, "token cannot contain whitespace")
		}
		if strings.Count(trimmed, ".") != 2 {
			errs = append(errs, "malformed token structure")
		}
	}
	if len(errs) > 0 {
		return false, 0, "", errors.New(strings.Join(errs, "; "))
	}

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
