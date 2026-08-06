package handlers

import (
	"errors"
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterAdminHandler(c *gin.Context) {
	var admin models.Admin

	if !BindJSON(c, &admin) {
		return
	}

	if err := models.Validate.Struct(admin); err != nil {
		errMap := models.ValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errMap})
		return
	}

	hashedPassword, err := HashPassword(admin.Password)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	admin.Password = hashedPassword

	id, err := h.service.RegisterAdmin(c.Request.Context(), admin)
	if err != nil {
		if errors.Is(err, services.ErrAdminEmailExists) {
			RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.TokenClient.GenerateToken(c.Request.Context(), id, admin.Email)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "admin created, but unable to generate token")
		return
	}

	RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      id,
		"token":   token,
		"message": "Admin created successfully",
	})
}
