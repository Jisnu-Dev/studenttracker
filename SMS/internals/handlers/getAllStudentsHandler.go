package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllStudentsHandler(c *gin.Context) {
	students, err := h.service.GetAllStudents(c.Request.Context())
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(c, http.StatusOK, students)
}
