package handlers

import (
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteStudentHandler(c *gin.Context) {
	id, ok := utils.ParseParamID(c, "id")
	if !ok {
		return
	}

	err := h.service.DeleteStudent(id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "Student deleted successfully",
	})
}
