package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

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

func checkError(t *testing.T, rr *httptest.ResponseRecorder, expected any) {
	t.Helper()
	switch exp := expected.(type) {
	case string:
		var resp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode error response: %v (body: %s)", err, rr.Body.String())
		}
		if resp.Error != exp {
			t.Errorf("expected error %q, got %q", exp, resp.Error)
		}
	case map[string]string:
		var resp struct {
			Errors map[string]string `json:"errors"`
		}
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode validation errors response: %v (body: %s)", err, rr.Body.String())
		}
		if !reflect.DeepEqual(resp.Errors, exp) {
			t.Errorf("expected validation errors %+v, got %+v", exp, resp.Errors)
		}
	default:
		t.Fatalf("unsupported expectedError type: %T", expected)
	}
}
