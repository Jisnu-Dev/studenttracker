package handlers

import (
	"errors"
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateStudentHandler(c *gin.Context) {
	id, ok := ParseParamID(c, "id")
	if !ok {
		return
	}

	var student models.Student
	if !BindJSON(c, &student) {
		return
	}

	if err := models.Validate.Struct(student); err != nil {
		errMap := models.ValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errMap})
		return
	}

	if err := h.service.UpdateStudent(c.Request.Context(), id, student); err != nil {
		if errors.Is(err, services.ErrStudentEmailExists) {
			RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(c, http.StatusOK, gin.H{"message": "Student updated successfully"})
}
