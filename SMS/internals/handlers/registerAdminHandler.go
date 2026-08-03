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

func (h *Handler) RegisterAdminHandler(c *gin.Context) {
	var admin models.Admin

	if !utils.BindJSON(c, &admin) {
		return
	}

	if err := validation.ValidateAdminRegister(&admin); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := utils.HashPassword(admin.Password)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	admin.Password = hashedPassword

	id, err := h.service.RegisterAdmin(admin)
	if err != nil {
		if errors.Is(err, services.ErrAdminEmailExists) {
			utils.RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.TokenClient.GenerateToken(c.Request.Context(), id, admin.Email)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "admin created, but unable to generate token")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      id,
		"token":   token,
		"message": "Admin created successfully",
	})
}
