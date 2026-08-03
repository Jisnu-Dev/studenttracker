package handlers_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	mockUtils "github.com/Jisnu-Dev/studenttracker/internals/mocks/utils"
	serviceErrors "github.com/Jisnu-Dev/studenttracker/internals/services/errors"
	"github.com/gin-gonic/gin"
)

func TestPatchStudentHandler(t *testing.T) {
	tests := []struct {
		name               string
		paramID            string
		body               string
		mockErr            mocks.MockOpError
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name:               "success - valid id and partial body (single field) patches student and returns 200",
			paramID:            "1",
			body:               `{"name":"test name"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "Student updated successfully",
			},
		},
		{
			name:               "success - valid id and missing name body patches student and returns 200",
			paramID:            "1",
			body:               `{"email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "Student updated successfully",
			},
		},
		{
			name:               "success - valid id and missing email body patches student and returns 200",
			paramID:            "1",
			body:               `{"name":"test name","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "Student updated successfully",
			},
		},
		{
			name:               "success - valid id and missing department body patches student and returns 200",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","semester":3,"age":20}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "Student updated successfully",
			},
		},
		{
			name:               "success - valid id and missing semester body patches student and returns 200",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","age":20}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "Student updated successfully",
			},
		},
		{
			name:               "success - valid id and missing age body patches student and returns 200",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","semester":3}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "Student updated successfully",
			},
		},
		{
			name:               "bad request - non-numeric id param returns 400",
			paramID:            "abc",
			body:               `{"name":"test name"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "invalid id parameter",
			},
		},
		{
			name:               "bad request - malformed JSON body returns 400",
			paramID:            "1",
			body:               `{`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "invalid request payload",
			},
		},
		{
			name:               "bad request - empty body with no fields returns 400",
			paramID:            "1",
			body:               `{}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "request: at least one field must be provided for update",
			},
		},
		{
			name:               "bad request - invalid email in patch fails validation",
			paramID:            "1",
			body:               `{"email":"not-an-email"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "email: email must contain '@'",
			},
		},
		{
			name:               "bad request - invalid semester in patch fails validation",
			paramID:            "1",
			body:               `{"semester":0}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "semester: semester must be between 1 and 8",
			},
		},
		{
			name:               "not found - student does not exist returns 404",
			paramID:            "99",
			body:               `{"name":"test name"}`,
			mockErr:            mocks.OpNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"error": "student not found",
			},
		},
		{
			name:               "conflict - duplicate email returns 409",
			paramID:            "1",
			body:               `{"email":"taken@example.com"}`,
			mockErr:            mocks.OpEmailExists,
			expectedStatusCode: http.StatusConflict,
			expectedResponse: map[string]interface{}{
				"error": "student with this email already exists",
			},
		},
		{
			name:               "bad request - no fields to update returns 400",
			paramID:            "1",
			body:               `{"name":"test name"}`,
			mockErr:            mocks.OpNoFieldsToUpdate,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": serviceErrors.ErrNoFieldsToUpdate.Error(),
			},
		},
		{
			name:               "internal server error - service failure returns 500",
			paramID:            "1",
			body:               `{"name":"test name"}`,
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "unable to patch student",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mockUtils.SetUpGinTestWithParams(
				http.MethodPatch, "/students/"+tt.paramID, tt.body,
				gin.Params{{Key: "id", Value: tt.paramID}},
			)

			svc := &mocks.MockService{PatchStudentError: tt.mockErr}
			handler := newHandler(svc, nil)
			handler.PatchStudentHandler(c)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedStatusCode, w.Code, w.Body.String())
			}

			if tt.expectedResponse != nil {
				var actualResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
				if err != nil {
					t.Fatalf("failed to parse response body as JSON: %v, body was: %s", err, w.Body.String())
				}

				if !reflect.DeepEqual(tt.expectedResponse, actualResponse) {
					t.Errorf("expected response %v, got %v", tt.expectedResponse, actualResponse)
				}
			}
		})
	}
}
