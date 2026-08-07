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

func TestUpdateStudentHandler(t *testing.T) {
	tests := []struct {
		name          string
		paramID       string
		body          string
		mockErr       mocks.MockOpError
		expectedCode  int
		expectedError any
	}{
		{
			name:         "update student successful",
			paramID:      "1",
			body:         `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			expectedCode: http.StatusOK,
		},
		{
			name:          "update student fails due to non-numeric id",
			paramID:       "abc",
			body:          `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidIDParam.Error(),
		},
		{
			name:          "update student fails due to invalid json",
			paramID:       "1",
			body:          `{`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidPayload.Error(),
		},
		{
			name:         "update student fails due to missing name",
			paramID:      "1",
			body:         `{"email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"name": models.ErrNameRequired.Error(),
			},
		},
		{
			name:         "update student fails due to missing email",
			paramID:      "1",
			body:         `{"name":"test name","department":"CSE","semester":3,"age":20}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email": models.ErrEmailRequired.Error(),
			},
		},
		{
			name:         "update student fails due to missing department",
			paramID:      "1",
			body:         `{"name":"test name","email":"test@example.com","semester":3,"age":20}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"department": models.ErrDepartmentRequired.Error(),
			},
		},
		{
			name:         "update student fails due to missing semester",
			paramID:      "1",
			body:         `{"name":"test name","email":"test@example.com","department":"CSE","age":20}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"semester": models.ErrSemesterInvalid.Error(),
			},
		},
		{
			name:         "update student fails due to missing age",
			paramID:      "1",
			body:         `{"name":"test name","email":"test@example.com","department":"CSE","semester":3}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"age": models.ErrAgeInvalid.Error(),
			},
		},
		{
			name:         "update student fails due to invalid department",
			paramID:      "1",
			body:         `{"name":"test name","email":"test@example.com","department":"INVALID","semester":3,"age":20}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"department": models.ErrInvalidDepartment.Error(),
			},
		},
		{
			name:         "update student successful when student not found (idempotent)",
			paramID:      "99",
			body:         `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			mockErr:      mocks.OpNotFound,
			expectedCode: http.StatusOK,
		},
		{
			name:          "update student fails because email already exists",
			paramID:       "1",
			body:          `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			mockErr:       mocks.OpEmailExists,
			expectedCode:  http.StatusConflict,
			expectedError: services.ErrStudentEmailExists.Error(),
		},
		{
			name:          "update student fails due to internal server error",
			paramID:       "1",
			body:          `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			mockErr:       mocks.OpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: services.ErrUpdateStudentFailed.Error(),
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
