package middlewares

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Jisnu-Dev/studenttracker/internals/grpcClient"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

var (
	ErrAuthHeaderRequired = errors.New("Authorization header is required")
	ErrAuthHeaderFormat   = errors.New("Authorization header format must be 'Bearer <token>'")
	ErrInvalidToken       = errors.New("Invalid or expired token")
)

func AuthMiddleware(tokenClient grpcClient.TokenClientInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrAuthHeaderRequired.Error(),
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrAuthHeaderFormat.Error(),
			})
			return
		}

		tokenString := parts[1]

		resp, err := tokenClient.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			st, _ := status.FromError(err)
			slog.Error("failed to validate token via gRPC",
				slog.String("middleware", "AuthMiddleware"),
				slog.String("grpcCode", st.Code().String()),
				slog.String("grpcMessage", st.Message()),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidToken.Error(),
			})
			return
		}

		if !resp.IsValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidToken.Error(),
			})
			return
		}

		c.Set("adminID", resp.AdminID)
		c.Set("adminEmail", resp.AdminEmail)

		c.Next()
	}
}
