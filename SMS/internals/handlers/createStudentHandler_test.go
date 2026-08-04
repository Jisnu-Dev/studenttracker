package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
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
		expectedBody       string
	}{
		{
			name:               "success - valid student body returns 200 with id and message",
			body:               mockUtils.ValidStudentJSON,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"id":1,"message":"student created successfully"}`,
		},
		{
			name:               "bad request - malformed JSON body returns 400",
			body:               `{`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid request payload"}`,
		},
		{
			name:               "bad request - missing name fails validation",
			body:               `{"email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"name: name is required"}`,
		},
		{
			name:               "bad request - missing email fails validation",
			body:               `{"name":"test name","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"email: email is required"}`,
		},
		{
			name:               "bad request - email missing @ symbol fails validation",
			body:               `{"name":"test name","email":"not-an-email","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"email: email must contain '@'"}`,
		},
		{
			name:               "bad request - structurally invalid email fails ParseAddress check",
			body:               `{"name":"test name","email":"invalid@","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"email: email is invalid"}`,
		},
		{
			name:               "bad request - invalid department value fails validation",
			body:               `{"name":"test name","email":"test@example.com","department":"INVALID","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"department: invalid department 'INVALID'; allowed values: CSE, IT, ECE, EEE, MECH, CIVIL, AIDS, AIML"}`,
		},
		{
			name:               "bad request - semester below minimum (0) fails validation",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","semester":0,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"semester: semester must be between 1 and 8"}`,
		},
		{
			name:               "bad request - semester above maximum (9) fails validation",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","semester":9,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"semester: semester must be between 1 and 8"}`,
		},
		{
			name:               "bad request - age below minimum (17) fails validation",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":17}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"age: age must be between 18 and 60"}`,
		},
		{
			name:               "bad request - age above maximum (61) fails validation",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":61}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"age: age must be between 18 and 60"}`,
		},
		{
			name:               "conflict - duplicate student email returns 409",
			body:               mockUtils.ValidStudentJSON,
			mockErr:            mocks.OpEmailExists,
			expectedStatusCode: http.StatusConflict,
			expectedBody:       `{"error":"student with this email already exists"}`,
		},
		{
			name:               "internal server error - service failure returns 500",
			body:               mockUtils.ValidStudentJSON,
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"unable to create student"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{CreateStudentError: tt.mockErr}
			_, mux := setupMockHandler(svc, nil)

			req := httptest.NewRequest(
				http.MethodPost,
				"/students",
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
