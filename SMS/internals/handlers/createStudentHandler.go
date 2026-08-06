package handlers

import (
	"errors"
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateStudentHandler(c *gin.Context) {
	var student models.Student

	if !BindJSON(c, &student) {
		return
	}

	if err := models.Validate.Struct(student); err != nil {
		errMap := models.ValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errMap})
		return
	}

	id, err := h.service.CreateStudent(c.Request.Context(), student)
	if err != nil {
		if errors.Is(err, services.ErrStudentEmailExists) {
			RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "student created successfully",
	})
}
