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

func (h *Handler) LoginAdminHandler(c *gin.Context) {
	var req models.LoginRequest

	if !utils.BindJSON(c, &req) {
		return
	}

	if err := validation.ValidateAdminLogin(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	admin, err := h.service.GetAdminByEmail(req.Email)
	if err != nil {
		if errors.Is(err, services.ErrAdminNotFound) {
			utils.RespondWithError(c, http.StatusUnauthorized, services.ErrInvalidCredentials.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !utils.CheckPasswordHash(req.Password, admin.Password) {
		utils.RespondWithError(c, http.StatusUnauthorized, services.ErrInvalidCredentials.Error())
		return
	}

	token, err := h.TokenClient.GenerateToken(c.Request.Context(), admin.ID, admin.Email)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "unable to login admin")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      admin.ID,
		"token":   token,
		"message": "Login successful",
	})
}
