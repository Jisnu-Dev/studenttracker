package handlers

import (
	"errors"
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/services"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetStudentByIDHandler(c *gin.Context) {
	id, ok := ParseParamID(c, "id")
	if !ok {
		return
	}

	student, err := h.service.GetStudentByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrStudentNotFound) {
			RespondWithJSON(c, http.StatusOK, gin.H{})
			return
		}
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(c, http.StatusOK, gin.H{"student": student})
}
