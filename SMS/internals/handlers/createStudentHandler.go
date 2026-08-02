package handlers

import (
	"net/http"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateStudentHandler(c *gin.Context) {
	var student models.Student

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//validate

	//call service
	id, err := h.service.CreateStudent(student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "student created successfully",
	})
}
