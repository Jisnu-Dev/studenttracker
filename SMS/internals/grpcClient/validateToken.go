package grpcClient

import (
	"context"
	"log/slog"
	"time"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
	"google.golang.org/grpc/status"
)

func (c *TokenClient) ValidateToken(ctx context.Context, token string) (bool, int64, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.ValidateToken(ctx, &tokenpb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("gRPC call failed",
			slog.String("function", "ValidateToken"),
			slog.String("grpcCode", st.Code().String()),
			slog.String("grpcMessage", st.Message()),
		)
		return false, 0, "", ErrValidateTokenFailed
	}

	return resp.GetIsValid(), resp.GetAdminID(), resp.GetAdminEmail(), nil
}
