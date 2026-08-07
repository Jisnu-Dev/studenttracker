package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetStudentByIDHandler(c *gin.Context) {
	id, ok := ParseParamID(c, "id")
	if !ok {
		return
	}

	student, err := h.service.GetStudentByID(c.Request.Context(), id)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if student.ID == 0 {
		RespondWithJSON(c, http.StatusOK, []interface{}{})
		return
	}

	RespondWithJSON(c, http.StatusOK, gin.H{"student": student})
}
