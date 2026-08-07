package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

func TestCreateStudentHandler(t *testing.T) {
	tests := []struct {
		name             string
		body             string
		mockErr       mocks.MockOpError
		expectedCode  int
		expectedBody  string
		expectedError any
	}{
		{
			name:         "student creation successful",
			body:         `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"message":"student created successfully"}`,
		},
		{
			name:          "student creation fails due to invalid json",
			body:          `{`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidPayload.Error(),
		},
		{
			name:         "student creation fails due to missing name",
			body:         `{"email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"name": models.ErrNameRequired.Error(),
			},
		},
		{
			name:         "student creation fails due to missing email",
			body:         `{"name":"test name","department":"CSE","semester":3,"age":20}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email": models.ErrEmailRequired.Error(),
			},
		},
		{
			name:         "student creation fails due to invalid email format",
			body:         `{"name":"test name","email":"not-an-email","department":"CSE","semester":3,"age":20}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email": models.ErrInvalidEmail.Error(),
			},
		},
		{
			name:         "student creation fails due to structurally invalid email",
			body:         `{"name":"test name","email":"invalid@","department":"CSE","semester":3,"age":20}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email": models.ErrInvalidEmail.Error(),
			},
		},
		{
			name:         "student creation fails due to invalid department",
			body:         `{"name":"test name","email":"test@example.com","department":"INVALID","semester":3,"age":20}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"department": models.ErrInvalidDepartment.Error(),
			},
		},
		{
			name:         "student creation fails due to semester below minimum (0)",
			body:         `{"name":"test name","email":"test@example.com","department":"CSE","semester":0,"age":20}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"semester": models.ErrSemesterInvalid.Error(),
			},
		},
		{
			name:         "student creation fails due to semester above maximum (9)",
			body:         `{"name":"test name","email":"test@example.com","department":"CSE","semester":9,"age":20}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"semester": models.ErrSemesterInvalid.Error(),
			},
		},
		{
			name:         "student creation fails due to age below minimum (17)",
			body:         `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":17}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"age": models.ErrAgeInvalid.Error(),
			},
		},
		{
			name:         "student creation fails due to age above maximum (61)",
			body:         `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":61}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"age": models.ErrAgeInvalid.Error(),
			},
		},
		{
			name:          "student creation fails because email already exists",
			body:          `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			mockErr:       mocks.OpEmailExists,
			expectedCode:  http.StatusConflict,
			expectedError: services.ErrStudentEmailExists.Error(),
		},
		{
			name:          "student creation fails due to internal server error",
			body:          `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":20}`,
			mockErr:       mocks.OpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: services.ErrCreateStudentFailed.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{
				CreateStudentError: tt.mockErr,
			}
			_, mux := setupMockHandler(svc, nil)

			req := httptest.NewRequest(
				http.MethodPost,
				"/students",
				bytes.NewBufferString(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedCode, rr.Code, rr.Body.String())
			}

			if tt.expectedBody != "" {
				got := rr.Body.String()
				if got[len(got)-1] == '\n' {
					got = got[:len(got)-1]
				}
				if got != tt.expectedBody {
					t.Errorf("expected body %q, got %q", tt.expectedBody, got)
				}
			}

			if tt.expectedError != nil {
				checkError(t, rr, tt.expectedError)
			}
		})
	}
}
