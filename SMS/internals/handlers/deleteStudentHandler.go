package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteStudentHandler(c *gin.Context) {
	id, ok := ParseParamID(c, "id")
	if !ok {
		return
	}

	if err := h.service.DeleteStudent(c.Request.Context(), id); err != nil {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "Student deleted successfully",
	})
}
