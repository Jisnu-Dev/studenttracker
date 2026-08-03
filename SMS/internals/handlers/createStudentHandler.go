package handlers

import (
	"errors"
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers/utils"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	services "github.com/Jisnu-Dev/studenttracker/internals/services/errors"
	"github.com/Jisnu-Dev/studenttracker/internals/validation"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateStudentHandler(c *gin.Context) {
	var student models.Student

	if !utils.BindJSON(c, &student) {
		return
	}

	if err := validation.ValidateStudent(&student); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateStudent(student)
	if err != nil {
		if errors.Is(err, services.ErrStudentEmailExists) {
			utils.RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "student created successfully",
	})
}
