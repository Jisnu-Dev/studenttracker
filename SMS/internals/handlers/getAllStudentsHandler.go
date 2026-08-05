package handlers

import (
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllStudentsHandler(c *gin.Context) {
	students, err := h.service.GetAllStudents()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if len(students) == 0 {
		utils.RespondWithJSON(c, http.StatusOK, gin.H{})
		return
	}
	utils.RespondWithJSON(c, http.StatusOK, students)
}
