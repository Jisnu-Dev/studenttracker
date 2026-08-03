package handlers

import (
	"errors"
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers/utils"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	serviceErrors "github.com/Jisnu-Dev/studenttracker/internals/services/errors"
	"github.com/Jisnu-Dev/studenttracker/internals/validation"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateStudentHandler(c *gin.Context) {
	id, ok := utils.ParseParamID(c, "id")
	if !ok {
		return
	}

	var student models.Student
	if !utils.BindJSON(c, &student) {
		return
	}

	if err := validation.ValidateStudent(&student); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateStudent(id, student); err != nil {
		if errors.Is(err, serviceErrors.ErrStudentNotFound) {
			utils.RespondWithError(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, serviceErrors.ErrStudentEmailExists) {
			utils.RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "Student updated successfully"})
}
