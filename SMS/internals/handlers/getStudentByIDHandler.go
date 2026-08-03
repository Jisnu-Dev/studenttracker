package handlers

import (
	"errors"
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers/utils"
	serviceErrors "github.com/Jisnu-Dev/studenttracker/internals/services/errors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetStudentByIDHandler(c *gin.Context) {
	id, ok := utils.ParseParamID(c, "id")
	if !ok {
		return
	}

	student, err := h.service.GetStudentByID(id)
	if err != nil {
		if errors.Is(err, serviceErrors.ErrStudentNotFound) {
			utils.RespondWithError(c, http.StatusNotFound, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"student": student})
}
