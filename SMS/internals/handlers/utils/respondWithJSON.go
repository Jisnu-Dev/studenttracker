package utils

import "github.com/gin-gonic/gin"

func RespondWithJSON(c *gin.Context, statusCode int, payload any) {
	c.JSON(statusCode, payload)
}
