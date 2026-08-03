package middlewares

import (
	"net/http"
	"strings"

	"github.com/Jisnu-Dev/studenttracker/internals/clients"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenClient clients.TokenClientInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header format must be 'Bearer <token>'",
			})
			return
		}

		tokenString := parts[1]

		isValid, adminID, adminEmail, err := tokenClient.ValidateToken(c.Request.Context(), tokenString)
		if err != nil || !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		c.Set("adminID", adminID)
		c.Set("adminEmail", adminEmail)

		c.Next()
	}
}
