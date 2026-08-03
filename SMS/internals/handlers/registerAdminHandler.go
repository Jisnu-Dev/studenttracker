package handlers

import (
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers/utils"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
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
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// call token service to generate token
	token, err := h.TokenClient.GenerateToken(c.Request.Context(), id, admin.Email)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Admin created, but token generation failed: "+err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      id,
		"token":   token,
		"message": "Admin created successfully",
	})
}
