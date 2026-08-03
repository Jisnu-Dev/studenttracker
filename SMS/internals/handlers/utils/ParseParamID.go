package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseParamID(c *gin.Context, paramName string) (int64, bool) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid "+paramName+" parameter: "+err.Error())
		return 0, false
	}
	return id, true
}
