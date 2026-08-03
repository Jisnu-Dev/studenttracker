package handlers_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	mockUtils "github.com/Jisnu-Dev/studenttracker/internals/mocks/utils"
)

func TestCreateStudentHandler(t *testing.T) {
	tests := []struct {
		name               string
		body               string
		mockErr            mocks.MockOpError
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name:               "success - valid student body returns 200 with id and message",
			body:               mockUtils.ValidStudentJSON,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "student created successfully",
			},
		},
		{
			name:               "bad request - malformed JSON body returns 400",
			body:               `{`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "invalid request payload",
			},
		},
		{
			name:               "bad request - missing name fails validation",
			body:               `{"email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "name: name is required",
			},
		},
		{
			name:               "bad request - missing email fails validation",
			body:               `{"name":"test name","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "email: email is required",
			},
		},
		{
			name:               "bad request - email missing @ symbol fails validation",
			body:               `{"name":"test name","email":"not-an-email","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "email: email must contain '@'",
			},
		},
		{
			name:               "bad request - structurally invalid email fails ParseAddress check",
			body:               `{"name":"test name","email":"invalid@","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "email: email is invalid",
			},
		},
		{
			name:               "bad request - invalid department value fails validation",
			body:               `{"name":"test name","email":"test@example.com","department":"INVALID","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "department: invalid department 'INVALID'; allowed values: CSE, IT, ECE, EEE, MECH, CIVIL, AIDS, AIML",
			},
		},
		{
			name:               "bad request - semester below minimum (0) fails validation",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","semester":0,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "semester: semester must be between 1 and 8",
			},
		},
		{
			name:               "bad request - semester above maximum (9) fails validation",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","semester":9,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "semester: semester must be between 1 and 8",
			},
		},
		{
			name:               "bad request - age below minimum (17) fails validation",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":17}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "age: age must be between 18 and 60",
			},
		},
		{
			name:               "bad request - age above maximum (61) fails validation",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":61}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "age: age must be between 18 and 60",
			},
		},
		{
			name:               "conflict - duplicate student email returns 409",
			body:               mockUtils.ValidStudentJSON,
			mockErr:            mocks.OpEmailExists,
			expectedStatusCode: http.StatusConflict,
			expectedResponse: map[string]interface{}{
				"error": "student with this email already exists",
			},
		},
		{
			name:               "internal server error - service failure returns 500",
			body:               mockUtils.ValidStudentJSON,
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "unable to create student",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mockUtils.SetUpGinTest(http.MethodPost, "/students", tt.body)

			svc := &mocks.MockService{CreateStudentError: tt.mockErr}
			handler := newHandler(svc, nil)
			handler.CreateStudentHandler(c)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedStatusCode, w.Code, w.Body.String())
			}

			if tt.expectedResponse != nil {
				var actualResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
				if err != nil {
					t.Fatalf("failed to parse response body as JSON: %v, body was: %s", err, w.Body.String())
				}

				for key, expectedValue := range tt.expectedResponse {
					actualValue, exists := actualResponse[key]
					if !exists {
						t.Errorf("expected key %q in response, but it was missing. Full response: %v", key, actualResponse)
						continue
					}
					if !reflect.DeepEqual(expectedValue, actualValue) {
						t.Errorf("expected %q to be %v, got %v", key, expectedValue, actualValue)
					}
				}
			}
		})
	}
}
