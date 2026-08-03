package handlers

import (
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers/utils"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateStudentHandler(c *gin.Context) {
	var student models.Student

	if !utils.BindJSON(c, &student) {
		return
	}

	//validate

	//call service
	id, err := h.service.CreateStudent(student)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "student created successfully",
	})
}
