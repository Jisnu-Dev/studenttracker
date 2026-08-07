package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/grpcClient"
	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
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
		name          string
		body          string
		mockErr       mocks.MockOpError
		generateErr   mocks.GrpcOpError
		mockToken     string
		mockAdmin     models.Admin
		expectedCode  int
		expectedError any
	}{
		{
			name:         "admin login successful",
			body:         `{"email":"admin@example.com","password":"Password123"}`,
			mockToken:    mocks.MockToken,
			expectedCode: http.StatusOK,
		},
		{
			name:          "admin login fails due to invalid json",
			body:          `{`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidPayload.Error(),
		},
		{
			name:         "admin login fails due to empty email",
			body:         `{"email":"","password":"Password123"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email": models.ErrEmailRequired.Error(),
			},
		},
		{
			name:         "admin login fails due to empty password",
			body:         `{"email":"admin@example.com","password":""}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"password": models.ErrPasswordRequired.Error(),
			},
		},
		{
			name:         "admin login fails due to empty email and password",
			body:         `{"email":"","password":""}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email":    models.ErrEmailRequired.Error(),
				"password": models.ErrPasswordRequired.Error(),
			},
		},
		{
			name:          "admin login fails due to invalid email format",
			body:          `{"email":"invalid-email","password":"Password123"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: map[string]string{"email": models.ErrInvalidEmail.Error()},
		},
		{
			name:          "admin login fails because email does not exist",
			body:          `{"email":"unknown@example.com","password":"Password123"}`,
			mockErr:       mocks.OpNotFound,
			expectedCode:  http.StatusUnauthorized,
			expectedError: services.ErrInvalidCredentials.Error(),
		},
		{
			name:          "email exists but internal server error",
			body:          `{"email":"admin@example.com","password":"Password123"}`,
			mockErr:       mocks.OpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlers.ErrUnableToLoginAdmin.Error(),
		},
		{
			name:          "email exists but password does not match",
			body:          `{"email":"admin@example.com","password":"WrongPassword1"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: services.ErrInvalidCredentials.Error(),
		},
		{
			name:          "email exists but TMS fails to generate token",
			body:          `{"email":"admin@example.com","password":"Password123"}`,
			generateErr:   mocks.GrpcOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlers.ErrUnableToLoginAdmin.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
				GetAdminByEmailError: tt.mockErr,
				Admin:                adminToReturn,
			}
			tc := grpcClient.NewTokenClientForTest(&mocks.MockTokenServiceClient{
				Token:       tt.mockToken,
				GenerateErr: tt.generateErr,
			})

			_, mux := setupMockHandler(svc, tc)

			req := httptest.NewRequest(
				http.MethodPost,
				"/login",
				bytes.NewBufferString(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected %d, got %d", tt.expectedCode, rr.Code)
			}

			if tt.expectedCode == http.StatusOK {
				resp := make(map[string]any)
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if resp["token"] == "" || resp["token"] == nil {
					t.Errorf("expected token in response, got empty")
				}
			}

			if tt.expectedError != nil {
				checkError(t, rr, tt.expectedError)
			}
		})
	}
}
