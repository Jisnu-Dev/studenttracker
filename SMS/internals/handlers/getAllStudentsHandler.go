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
	if len(students) == 0 {
		RespondWithJSON(c, http.StatusOK, gin.H{})
		return
	}
	RespondWithJSON(c, http.StatusOK, students)
}
