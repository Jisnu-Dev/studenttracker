package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/grpcClient"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
)

func TestRegisterAdminHandler(t *testing.T) {
	tests := []struct {
		name               string
		body               string
		mockErr            mocks.MockOpError
		generateTokenErr   mocks.GrpcOpError
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "success - valid admin body returns 200 with id and token",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234!"}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"id":1,"message":"Admin created successfully","token":"mock.jwt.token"}`,
		},
		{
			name:               "bad request - malformed JSON body returns 400",
			body:               `{`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid request payload"}`,
		},
		{
			name:               "bad request - empty name fails validation",
			body:               `{"name":"","email":"john.admin@example.com","password":"Admin1234!"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":{"name":"name is required"}}`,
		},
		{
			name:               "bad request - empty email fails validation",
			body:               `{"name":"John Admin","email":"","password":"Admin1234!"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":{"email":"email is required"}}`,
		},
		{
			name:               "bad request - empty password fails validation",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":""}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":{"password":"password is required"}}`,
		},
		{
			name:               "bad request - password too short fails validation",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"ab1"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":{"password":"password must be at least 8 characters"}}`,
		},
		{
			name:               "bad request - password with no digit fails validation",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"OnlyLetters!"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":{"password":"password must contain uppercase, lowercase, digit, and special character"}}`,
		},
		{
			name:               "conflict - duplicate admin email returns 409",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234!"}`,
			mockErr:            mocks.OpEmailExists,
			expectedStatusCode: http.StatusConflict,
			expectedBody:       `{"error":"admin with this email already exists"}`,
		},
		{
			name:               "internal server error - service failure returns 500",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234!"}`,
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"unable to register admin"}`,
		},
		{
			name:               "internal server error - password hashing failure returns 500",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1!Aa` + strings.Repeat("✨", 30) + `"}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"Failed to hash password"}`,
		},
		{
			name:               "internal server error - token generation failure after admin created returns 500",
			body:               `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234!"}`,
			generateTokenErr:   mocks.GrpcOpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"admin created, but unable to generate token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{RegisterAdminError: tt.mockErr}
			tc := grpcClient.NewTokenClientForTest(&mocks.MockTokenServiceClient{GenerateErr: tt.generateTokenErr})
			
			_, mux := setupMockHandler(svc, tc)

			req := httptest.NewRequest(
				http.MethodPost,
				"/register",
				bytes.NewBufferString(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")
			
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedStatusCode, rr.Code, rr.Body.String())
			}

			if got := strings.TrimSpace(rr.Body.String()); got != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, got)
			}
		})
	}
}
