package handlers_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	mockUtils "github.com/Jisnu-Dev/studenttracker/internals/mocks/utils"
	"github.com/gin-gonic/gin"
)

func TestUpdateStudentHandler(t *testing.T) {
	tests := []struct {
		name               string
		paramID            string
		body               string
		mockErr            mocks.MockOpError
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name:               "success - valid id and body updates student and returns 200",
			paramID:            "1",
			body:               mockUtils.ValidStudentJSON,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "Student updated successfully",
			},
		},
		{
			name:               "bad request - non-numeric id param returns 400",
			paramID:            "abc",
			body:               mockUtils.ValidStudentJSON,
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
			name:               "bad request - missing name fails validation",
			paramID:            "1",
			body:               `{"email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "name: name is required",
			},
		},
		{
			name:               "bad request - missing email fails validation",
			paramID:            "1",
			body:               `{"name":"test name","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "email: email is required",
			},
		},
		{
			name:               "bad request - missing department fails validation",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "department: invalid department ''; allowed values: CSE, IT, ECE, EEE, MECH, CIVIL, AIDS, AIML",
			},
		},
		{
			name:               "bad request - missing semester fails validation",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "semester: semester must be between 1 and 8",
			},
		},
		{
			name:               "bad request - missing age fails validation",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","semester":3}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "age: age must be between 18 and 60",
			},
		},
		{
			name:               "bad request - invalid department fails validation",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","department":"INVALID","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "department: invalid department 'INVALID'; allowed values: CSE, IT, ECE, EEE, MECH, CIVIL, AIDS, AIML",
			},
		},
		{
			name:               "not found - student does not exist returns 404",
			paramID:            "99",
			body:               mockUtils.ValidStudentJSON,
			mockErr:            mocks.OpNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"error": "student not found",
			},
		},
		{
			name:               "conflict - duplicate email returns 409",
			paramID:            "1",
			body:               mockUtils.ValidStudentJSON,
			mockErr:            mocks.OpEmailExists,
			expectedStatusCode: http.StatusConflict,
			expectedResponse: map[string]interface{}{
				"error": "student with this email already exists",
			},
		},
		{
			name:               "internal server error - service failure returns 500",
			paramID:            "1",
			body:               mockUtils.ValidStudentJSON,
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "unable to update student",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mockUtils.SetUpGinTestWithParams(
				http.MethodPut, "/students/"+tt.paramID, tt.body,
				gin.Params{{Key: "id", Value: tt.paramID}},
			)

			svc := &mocks.MockService{UpdateStudentError: tt.mockErr}
			handler := newHandler(svc, nil)
			handler.UpdateStudentHandler(c)

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
