package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/grpcClient"
	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

func TestRegisterAdminHandler(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		mockErr       mocks.MockOpError
		generateErr   mocks.GrpcOpError
		mockToken     string
		registeredID  int64
		expectedCode  int
		expectedError any
	}{
		{
			name:         "admin registration successful",
			body:         `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234!"}`,
			mockToken:    mocks.MockToken,
			registeredID: 1,
			expectedCode: http.StatusOK,
		},
		{
			name:          "admin registration fails due to invalid json",
			body:          `{`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidPayload.Error(),
		},
		{
			name:         "admin registration fails due to empty name",
			body:         `{"name":"","email":"john.admin@example.com","password":"Admin1234!"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"name": models.ErrNameRequired.Error(),
			},
		},
		{
			name:         "admin registration fails due to empty email",
			body:         `{"name":"John Admin","email":"","password":"Admin1234!"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email": models.ErrEmailRequired.Error(),
			},
		},
		{
			name:         "admin registration fails due to empty password",
			body:         `{"name":"John Admin","email":"john.admin@example.com","password":""}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"password": models.ErrPasswordRequired.Error(),
			},
		},
		{
			name:         "admin registration fails due to password too short",
			body:         `{"name":"John Admin","email":"john.admin@example.com","password":"ab1"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"password": models.ErrPasswordTooShort.Error(),
			},
		},
		{
			name:         "admin registration fails due to password with no digit",
			body:         `{"name":"John Admin","email":"john.admin@example.com","password":"OnlyLetters!"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"password": models.ErrInvalidPassword.Error(),
			},
		},
		{
			name:          "admin registration fails because email already exists",
			body:          `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234!"}`,
			mockErr:       mocks.OpEmailExists,
			expectedCode:  http.StatusConflict,
			expectedError: services.ErrAdminEmailExists.Error(),
		},
		{
			name:          "admin registration fails due to internal server error",
			body:          `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234!"}`,
			mockErr:       mocks.OpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlers.ErrUnableToRegisterAdmin.Error(),
		},
		{
			name:          "admin registration fails due to password hashing failure",
			body:          `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1!Aa` + strings.Repeat("✨", 30) + `"}`,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlers.ErrUnableToRegisterAdmin.Error(),
		},
		{
			name:          "admin registration fails when TMS fails to generate token",
			body:          `{"name":"John Admin","email":"john.admin@example.com","password":"Admin1234!"}`,
			generateErr:   mocks.GrpcOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlers.ErrUnableToRegisterAdmin.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{
				RegisterAdminError: tt.mockErr,
				RegisteredAdminID:  tt.registeredID,
			}
			tc := grpcClient.NewTokenClientForTest(&mocks.MockTokenServiceClient{
				Token:       tt.mockToken,
				GenerateErr: tt.generateErr,
			})

			_, mux := setupMockHandler(svc, tc)

			req := httptest.NewRequest(
				http.MethodPost,
				"/register",
				bytes.NewBufferString(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedCode, rr.Code, rr.Body.String())
			}

			if tt.expectedCode == http.StatusOK {
				var resp struct {
					ID      int64  `json:"id"`
					Token   string `json:"token"`
					Message string `json:"message"`
				}
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if resp.Token == "" {
					t.Errorf("expected token in response, got empty")
				}
				if resp.Message != "Admin created successfully" {
					t.Errorf("expected message 'Admin created successfully', got %q", resp.Message)
				}
			}

			if tt.expectedError != nil {
				checkError(t, rr, tt.expectedError)
			}
		})
	}
}
