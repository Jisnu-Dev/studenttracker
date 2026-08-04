package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
)

func TestPatchStudentHandler(t *testing.T) {
	tests := []struct {
		name               string
		paramID            string
		body               string
		mockErr            mocks.MockOpError
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "success - valid id and partial body (single field) patches student and returns 200",
			paramID:            "1",
			body:               `{"name":"test name"}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"Student updated successfully"}`,
		},
		{
			name:               "success - valid id and missing name body patches student and returns 200",
			paramID:            "1",
			body:               `{"email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"Student updated successfully"}`,
		},
		{
			name:               "success - valid id and missing email body patches student and returns 200",
			paramID:            "1",
			body:               `{"name":"test name","department":"CSE","semester":3,"age":20}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"Student updated successfully"}`,
		},
		{
			name:               "success - valid id and missing department body patches student and returns 200",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","semester":3,"age":20}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"Student updated successfully"}`,
		},
		{
			name:               "success - valid id and missing semester body patches student and returns 200",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","age":20}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"Student updated successfully"}`,
		},
		{
			name:               "success - valid id and missing age body patches student and returns 200",
			paramID:            "1",
			body:               `{"name":"test name","email":"test@example.com","department":"CSE","semester":3}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"Student updated successfully"}`,
		},
		{
			name:               "bad request - non-numeric id param returns 400",
			paramID:            "abc",
			body:               `{"name":"test name"}`,
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
			name:               "bad request - empty body with no fields returns 400",
			paramID:            "1",
			body:               `{}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"request: at least one field must be provided for update"}`,
		},
		{
			name:               "bad request - invalid email in patch fails validation",
			paramID:            "1",
			body:               `{"email":"not-an-email"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"email: email must contain '@'"}`,
		},
		{
			name:               "bad request - invalid semester in patch fails validation",
			paramID:            "1",
			body:               `{"semester":0}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"semester: semester must be between 1 and 8"}`,
		},
		{
			name:               "not found - student does not exist returns 404",
			paramID:            "99",
			body:               `{"name":"test name"}`,
			mockErr:            mocks.OpNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `{"error":"student not found"}`,
		},
		{
			name:               "conflict - duplicate email returns 409",
			paramID:            "1",
			body:               `{"email":"taken@example.com"}`,
			mockErr:            mocks.OpEmailExists,
			expectedStatusCode: http.StatusConflict,
			expectedBody:       `{"error":"student with this email already exists"}`,
		},
		{
			name:               "bad request - no fields to update returns 400",
			paramID:            "1",
			body:               `{"name":"test name"}`,
			mockErr:            mocks.OpNoFieldsToUpdate,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"at least one field must be provided for update"}`,
		},
		{
			name:               "internal server error - service failure returns 500",
			paramID:            "1",
			body:               `{"name":"test name"}`,
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"unable to patch student"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{PatchStudentError: tt.mockErr}
			_, mux := setupMockHandler(svc, nil)

			req := httptest.NewRequest(
				http.MethodPatch,
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
