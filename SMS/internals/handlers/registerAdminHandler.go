package handlers

import (
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) RegisterAdminHandler(c *gin.Context) {
	var admin models.Admin

	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(admin.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin.Password = string(hashedPassword)

	id, err := h.service.RegisterAdmin(admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// call token service to generate token
	token, err := h.TokenClient.GenerateToken(c.Request.Context(), id, admin.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Admin created, but token generation failed: ": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"token":   token,
		"message": "Admin created successfully",
	})
}
