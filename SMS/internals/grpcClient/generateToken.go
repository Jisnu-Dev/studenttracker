package grpcClient

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
	"github.com/Jisnu-Dev/studenttracker/internals/validation"
)

func (c *TokenClient) GenerateToken(ctx context.Context, adminID int64, adminEmail string) (string, error) {
	var errs []string
	if adminID <= 0 {
		errs = append(errs, "admin_id is required and must be greater than 0")
	}
	if err := validation.ValidateEmail(adminEmail); err != nil {
		errs = append(errs, "admin_email: "+err.Error())
	}
	if len(errs) > 0 {
		return "", errors.New(strings.Join(errs, "; "))
	}

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
