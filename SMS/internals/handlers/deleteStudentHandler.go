package handlers

import (
	"errors"
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/services"
	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteStudentHandler(c *gin.Context) {
	id, ok := ParseParamID(c, "id")
	if !ok {
		return
	}

	err := h.service.DeleteStudent(c.Request.Context(), id)
	if err != nil && !errors.Is(err, services.ErrStudentNotFound) {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "Student deleted successfully",
	})
}
