package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
)

func TestUpdateStudentHandler(t *testing.T) {
	tests := []struct {
		name               string
		paramID            string
		body               string
		mockErr            mocks.MockOpError
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "success - valid id and body updates student and returns 200",
			paramID:            "1",
			body:               validStudentJSON,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"Student updated successfully"}`,
		},
		{
			name:               "bad request - non-numeric id param returns 400",
			paramID:            "abc",
			body:               validStudentJSON,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid id parameter"}`,
		},
		{
			name:               "bad request - malformed JSON body returns 400",
			paramID:            "1",
			body:               `{`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid request payload"}`,
		},
		{
			name:               "bad request - missing name fails validation",
			paramID:            "1",
			body:               `{"email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":{"name":"name is required"}}`,
		},
		{
			name:               "bad request - missing email fails validation",
			paramID:            "1",
			body:               `{"name":"test name","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":{"email":"email is required"}}`,
		},
		{
			name:               "bad request - missing department fails validation",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":{"department":"department is required"}}`,
		},
		{
			name:               "bad request - missing semester fails validation",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":{"semester":"semester must be between 1 and 8"}}`,
		},
		{
			name:               "bad request - missing age fails validation",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","semester":3}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":{"age":"age must be between 18 and 60"}}`,
		},
		{
			name:               "bad request - invalid department fails validation",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","department":"INVALID","semester":3,"age":20}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":{"department":"invalid department"}}`,
		},
		{
			name:               "success (not found) - student does not exist but idempotent update returns 200",
			paramID:            "99",
			body:               validStudentJSON,
			mockErr:            mocks.OpNotFound,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"Student updated successfully"}`,
		},
		{
			name:               "conflict - duplicate email returns 409",
			paramID:            "1",
			body:               validStudentJSON,
			mockErr:            mocks.OpEmailExists,
			expectedStatusCode: http.StatusConflict,
			expectedBody:       `{"error":"student with this email already exists"}`,
		},
		{
			name:               "internal server error - service failure returns 500",
			paramID:            "1",
			body:               validStudentJSON,
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"unable to update student"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{UpdateStudentError: tt.mockErr}
			_, mux := setupMockHandler(svc, nil)

			req := httptest.NewRequest(
				http.MethodPut,
				"/students/"+tt.paramID,
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
