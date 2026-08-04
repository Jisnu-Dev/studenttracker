package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"golang.org/x/crypto/bcrypt"
)

// loginTestPassword is the plaintext password used across login handler tests.
const loginTestPassword = "Password123"

// hashedLoginTestPassword is the bcrypt hash of loginTestPassword, generated once at package init.
var hashedLoginTestPassword string

func init() {
	hash, err := bcrypt.GenerateFromPassword([]byte(loginTestPassword), bcrypt.MinCost)
	if err != nil {
		panic("failed to hash login test password: " + err.Error())
	}
	hashedLoginTestPassword = string(hash)
}

func TestLoginAdminHandler(t *testing.T) {
	tests := []struct {
		name               string
		body               string
		getAdminErr        mocks.MockOpError
		generateTokenErr   mocks.MockOpError
		mockAdmin          models.Admin
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "success - valid credentials returns 200 with token",
			body:               `{"email":"admin@example.com","password":"` + loginTestPassword + `"}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"id":1,"message":"Login successful","token":"mock.jwt.token"}`,
		},
		{
			name:               "bad request - malformed JSON body returns 400",
			body:               `{`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid request payload"}`,
		},
		{
			name:               "bad request - empty email fails validation",
			body:               `{"email":"","password":"Password123"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"email: email is required"}`,
		},
		{
			name:               "bad request - empty password fails validation",
			body:               `{"email":"admin@example.com","password":""}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"password: password is required"}`,
		},
		{
			name:               "unauthorized - admin not found returns 401 with invalid credentials message",
			body:               `{"email":"unknown@example.com","password":"Password123"}`,
			getAdminErr:        mocks.OpNotFound,
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"invalid email or password"}`,
		},
		{
			name:               "internal server error - get admin service failure returns 500",
			body:               `{"email":"admin@example.com","password":"Password123"}`,
			getAdminErr:        mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"unable to retrieve admin"}`,
		},
		{
			name:               "unauthorized - wrong password returns 401 with invalid credentials message",
			body:               `{"email":"admin@example.com","password":"WrongPassword1"}`,
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"invalid email or password"}`,
		},
		{
			name:               "internal server error - token generation failure returns 500",
			body:               `{"email":"admin@example.com","password":"` + loginTestPassword + `"}`,
			generateTokenErr:   mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"unable to login admin"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Seed the mock admin with a properly hashed password so CheckPasswordHash works.
			adminToReturn := models.Admin{
				ID:       1,
				Name:     "Test Admin",
				Email:    "admin@example.com",
				Password: hashedLoginTestPassword,
			}
			if tt.mockAdmin.Email != "" {
				adminToReturn = tt.mockAdmin
			}

			svc := &mocks.MockService{
				GetAdminByEmailError: tt.getAdminErr,
				Admin:                adminToReturn,
			}
			tc := &mocks.MockTokenClient{
				GenerateTokenError: tt.generateTokenErr,
			}
			
			_, mux := setupMockHandler(svc, tc)

			req := httptest.NewRequest(
				http.MethodPost,
				"/login",
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
