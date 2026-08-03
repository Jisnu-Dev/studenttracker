package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers/utils"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	mockUtils "github.com/Jisnu-Dev/studenttracker/internals/mocks/utils"
)

func TestRegisterAdminHandler(t *testing.T) {
	tests := []struct {
		name               string
		body               string
		mockErr            mocks.MockOpError
		generateTokenErr   mocks.MockOpError
		expectedStatusCode int
		expectedResponse   map[string]interface{}
		simulateHashErr    bool
	}{
		{
			name:               "success - valid admin body returns 200 with id and token",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"id":      float64(1), // Default mock returned ID
				"token":   "mock.jwt.token",
				"message": "Admin created successfully",
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
			name:               "bad request - empty name fails validation",
			body:               `{"name":"","email":"john.admin@example.com","password":"Admin1234"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "name: name is required",
			},
		},
		{
			name:               "bad request - empty email fails validation",
			body:               `{"name":"John Admin","email":"","password":"Admin1234"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "email: email is required",
			},
		},
		{
			name:               "bad request - empty password fails validation",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":""}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "password: password is required",
			},
		},
		{
			name:               "bad request - password too short fails validation",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"ab1"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "password: password must be at least 8 characters long",
			},
		},
		{
			name:               "bad request - password with no digit fails validation",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"OnlyLetters"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "password: password must contain at least one letter and one number",
			},
		},
		{
			name:               "conflict - duplicate admin email returns 409",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234"}`,
			mockErr:            mocks.OpEmailExists,
			expectedStatusCode: http.StatusConflict,
			expectedResponse: map[string]interface{}{
				"error": "admin with this email already exists",
			},
		},
		{
			name:               "internal server error - service failure returns 500",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234"}`,
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "unable to register admin",
			},
		},
		{
			name:               "internal server error - token generation failure after admin created returns 500",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234"}`,
			generateTokenErr:   mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "admin created, but unable to generate token",
			},
		},
		{
			name:               "internal server error - hash password failure returns 500",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234"}`,
			simulateHashErr:    true,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "Failed to hash password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.simulateHashErr {
				originalHash := utils.HashPassword
				utils.HashPassword = func(password string) (string, error) {
					return "", fmt.Errorf("mock hash error")
				}
				defer func() { utils.HashPassword = originalHash }()
			}

			c, w := mockUtils.SetUpGinTest(http.MethodPost, "/register", tt.body)

			svc := &mocks.MockService{RegisterAdminError: tt.mockErr}
			tc := &mocks.MockTokenClient{GenerateTokenError: tt.generateTokenErr}
			handler := newHandler(svc, tc)
			handler.RegisterAdminHandler(c)

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
