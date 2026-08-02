package handlers

import (
	"net/http"
	"strconv"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) PatchStudentHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var patchStudent models.PatchStudent
	if err := c.ShouldBindJSON(&patchStudent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.PatchStudent(idInt, patchStudent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student updated successfully"})
}
