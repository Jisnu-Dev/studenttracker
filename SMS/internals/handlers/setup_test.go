package handlers_test

import (
	"github.com/Jisnu-Dev/studenttracker/internals/grpcClient"
	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/gin-gonic/gin"
)

func setupMockHandler(svc *mocks.MockService, tc grpcClient.TokenClientInterface) (*handlers.Handler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := handlers.NewHandler(svc, tc)

	router.POST("/students", handler.CreateStudentHandler)
	router.DELETE("/students/:id", handler.DeleteStudentHandler)
	router.GET("/students", handler.GetAllStudentsHandler)
	router.GET("/students/:id", handler.GetStudentByIDHandler)
	router.PATCH("/students/:id", handler.PatchStudentHandler)
	router.PUT("/students/:id", handler.UpdateStudentHandler)
	router.POST("/login", handler.LoginAdminHandler)
	router.POST("/register", handler.RegisterAdminHandler)

	return handler, router
}
