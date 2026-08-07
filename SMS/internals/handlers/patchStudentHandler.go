package handlers

import (
	"errors"
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
	"github.com/gin-gonic/gin"
)

func (h *Handler) PatchStudentHandler(c *gin.Context) {
	id, ok := ParseParamID(c, "id")
	if !ok {
		return
	}

	var patchStudent models.PatchStudent
	if !BindJSON(c, &patchStudent) {
		return
	}

	if patchStudent.Name == nil && patchStudent.Email == nil && patchStudent.Department == nil && patchStudent.Semester == nil && patchStudent.Age == nil {
		RespondWithError(c, http.StatusBadRequest, models.ErrAtLeastOneField.Error())
		return
	}

	if err := models.Validate.Struct(patchStudent); err != nil {
		errMap := models.ValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errMap})
		return
	}

	if err := h.service.PatchStudent(c.Request.Context(), id, patchStudent); err != nil {
		if errors.Is(err, services.ErrStudentEmailExists) {
			RespondWithError(c, http.StatusConflict, err.Error())
			return
		}
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(c, http.StatusOK, gin.H{"message": "Student updated successfully"})
}
