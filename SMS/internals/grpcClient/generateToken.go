package grpcClient

import (
	"context"
	"log/slog"
	"time"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
	"google.golang.org/grpc/status"
)

func (c *TokenClient) GenerateToken(ctx context.Context, adminID int64, adminEmail string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.GenerateToken(ctx, &tokenpb.GenerateTokenRequest{
		AdminID:    adminID,
		AdminEmail: adminEmail,
	})
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("gRPC call failed",
			slog.String("function", "GenerateToken"),
			slog.String("adminEmail", adminEmail),
			slog.String("grpcCode", st.Code().String()),
			slog.String("grpcMessage", st.Message()),
		)
		return "", ErrGenerateTokenFailed
	}

	return resp.GetToken(), nil
}
