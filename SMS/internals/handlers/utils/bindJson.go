package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindJSON(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		RespondWithError(c, http.StatusBadRequest, "invalid request payload")
		return false
	}
	return true
}
