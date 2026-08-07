package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

func TestPatchStudentHandler(t *testing.T) {
	tests := []struct {
		name          string
		paramID       string
		body          string
		mockErr       mocks.MockOpError
		expectedCode  int
		expectedError any
	}{
		{
			name:         "patch student successful - single field",
			paramID:      "1",
			body:         `{"name":"test name"}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "patch student successful - missing name",
			paramID:      "1",
			body:         `{"email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "patch student successful - missing email",
			paramID:      "1",
			body:         `{"name":"test name","department":"CSE","semester":3,"age":20}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "patch student successful - missing department",
			paramID:      "1",
			body:         `{"name":"test name","email":"test@example.com","semester":3,"age":20}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "patch student successful - missing semester",
			paramID:      "1",
			body:         `{"name":"test name","email":"test@example.com","department":"CSE","age":20}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "patch student successful - missing age",
			paramID:      "1",
			body:         `{"name":"test name","email":"test@example.com","department":"CSE","semester":3}`,
			expectedCode: http.StatusOK,
		},
		{
			name:          "patch student fails due to non-numeric id",
			paramID:       "abc",
			body:          `{"name":"test name"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidIDParam.Error(),
		},
		{
			name:          "patch student fails due to invalid json",
			paramID:       "1",
			body:          `{`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidPayload.Error(),
		},
		{
			name:          "patch student fails due to empty body with no fields",
			paramID:       "1",
			body:          `{}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: models.ErrAtLeastOneField.Error(),
		},
		{
			name:         "patch student fails due to invalid email format",
			paramID:      "1",
			body:         `{"email":"not-an-email"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email": models.ErrInvalidEmail.Error(),
			},
		},
		{
			name:         "patch student fails due to invalid semester",
			paramID:      "1",
			body:         `{"semester":0}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"semester": models.ErrSemesterInvalid.Error(),
			},
		},
		{
			name:         "patch student successful when student not found (idempotent)",
			paramID:      "99",
			body:         `{"name":"test name"}`,
			mockErr:      mocks.OpNotFound,
			expectedCode: http.StatusOK,
		},
		{
			name:          "patch student fails because email already exists",
			paramID:       "1",
			body:          `{"email":"taken@example.com"}`,
			mockErr:       mocks.OpEmailExists,
			expectedCode:  http.StatusConflict,
			expectedError: services.ErrStudentEmailExists.Error(),
		},
		{
			name:          "patch student fails due to internal server error",
			paramID:       "1",
			body:          `{"name":"test name"}`,
			mockErr:       mocks.OpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: services.ErrPatchStudentFailed.Error(),
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

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedCode, rr.Code, rr.Body.String())
			}

			if tt.expectedCode == http.StatusOK {
				var resp struct {
					Message string `json:"message"`
				}
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if resp.Message != "Student updated successfully" {
					t.Errorf("expected message 'Student updated successfully', got %q", resp.Message)
				}
			}

			if tt.expectedError != nil {
				checkError(t, rr, tt.expectedError)
			}
		})
	}
}
