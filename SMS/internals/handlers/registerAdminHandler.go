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
		slog.Error("failed to hash admin password",
			slog.String("handler", "RegisterAdminHandler"),
			slog.String("adminEmail", admin.Email),
			slog.Any("error", err),
		)
		RespondWithError(c, http.StatusInternalServerError, ErrUnableToRegisterAdmin.Error())
		return
	}

	admin.Password = hashedPassword

	id, err := h.service.RegisterAdmin(c.Request.Context(), admin)
	if err != nil {
		if errors.Is(err, services.ErrAdminEmailExists) {
			RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		RespondWithError(c, http.StatusInternalServerError, ErrUnableToRegisterAdmin.Error())
		return
	}

	token, err := h.TokenClient.GenerateToken(c.Request.Context(), id, admin.Email)
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("failed to generate token for registered admin",
			slog.String("handler", "RegisterAdminHandler"),
			slog.Int64("adminID", id),
			slog.String("adminEmail", admin.Email),
			slog.String("grpcCode", st.Code().String()),
			slog.String("grpcMessage", st.Message()),
		)
		RespondWithError(c, http.StatusInternalServerError, ErrUnableToRegisterAdmin.Error())
		return
	}

	RespondWithJSON(c, http.StatusOK, gin.H{
		"id":      id,
		"token":   token,
		"message": "Admin created successfully",
	})
}
