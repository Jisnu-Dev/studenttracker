package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

func (h *Handler) LoginAdminHandler(c *gin.Context) {
	var req models.LoginRequest

	if !BindJSON(c, &req) {
		return
	}

	if err := models.Validate.Struct(req); err != nil {
		errMap := models.ValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errMap})
		return
	}

	admin, err := h.service.GetAdminByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, services.ErrAdminNotFound) {
			RespondWithError(c, http.StatusUnauthorized, services.ErrInvalidCredentials.Error())
			return
		}
		RespondWithError(c, http.StatusInternalServerError, ErrUnableToLoginAdmin.Error())
		return
	}

	if !CheckPasswordHash(req.Password, admin.Password) {
		RespondWithError(c, http.StatusUnauthorized, services.ErrInvalidCredentials.Error())
		return
	}

	token, err := h.TokenClient.GenerateToken(c.Request.Context(), admin.ID, admin.Email)
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("failed to generate token for login",
			slog.String("handler", "LoginAdminHandler"),
			slog.Int64("adminID", admin.ID),
			slog.String("adminEmail", admin.Email),
			slog.String("grpcCode", st.Code().String()),
			slog.String("grpcMessage", st.Message()),
		)
		RespondWithError(c, http.StatusInternalServerError, ErrUnableToLoginAdmin.Error())
		return
	}

	RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      admin.ID,
		"token":   token,
		"message": "Login successful",
	})
}
