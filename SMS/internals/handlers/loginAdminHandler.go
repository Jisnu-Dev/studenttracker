package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers/utils"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	serviceErrors "github.com/Jisnu-Dev/studenttracker/internals/services/errors"
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

	slog.Info("Login request: %v", slog.String("email", req.Email))
	admin, err := h.service.GetAdminByEmail(req.Email)
	if err != nil {
		if errors.Is(err, serviceErrors.ErrAdminNotFound) {
			utils.RespondWithError(c, http.StatusUnauthorized, serviceErrors.ErrInvalidCredentials.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !utils.CheckPasswordHash(req.Password, admin.Password) {
		utils.RespondWithError(c, http.StatusUnauthorized, serviceErrors.ErrInvalidCredentials.Error())
		return
	}

	token, err := h.TokenClient.GenerateToken(c.Request.Context(), admin.ID, admin.Email)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate token: "+err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      admin.ID,
		"token":   token,
		"message": "Login successful",
	})
}
