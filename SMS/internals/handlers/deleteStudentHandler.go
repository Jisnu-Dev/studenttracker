package handlers

import (
	"errors"
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers/utils"
	services "github.com/Jisnu-Dev/studenttracker/internals/services/errors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteStudentHandler(c *gin.Context) {
	id, ok := utils.ParseParamID(c, "id")
	if !ok {
		return
	}

	err := h.service.DeleteStudent(id)
	if err != nil && !errors.Is(err, services.ErrStudentNotFound) {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      id,
		"message": "Student deleted successfully",
	})
}
