package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

func TestGetAllStudentsHandler(t *testing.T) {
	tests := []struct {
		name          string
		mockErr       mocks.MockOpError
		returnEmpty   bool
		mockStudents  []models.Student
		expectedBody  string
		expectedCode  int
		expectedError any
	}{
		{
			name: "retrieve all students successful - single student",
			expectedBody: `[{"id":1,"name":"test","email":"test@example.com","department":"CSE","semester":3,"age":20,"createdAtUtc":"0001-01-01T00:00:00Z","updatedAtUtc":"0001-01-01T00:00:00Z"}]`,
			expectedCode: http.StatusOK,
		},
		{
			name: "retrieve all students successful - multiple students",
			mockStudents: []models.Student{
				{ID: 1, Name: "Alice", Email: "alice@example.com", Department: models.CSE, Semester: 5, Age: 21},
				{ID: 2, Name: "Bob", Email: "bob@example.com", Department: models.ECE, Semester: 3, Age: 20},
				{ID: 3, Name: "Charlie", Email: "charlie@example.com", Department: models.IT, Semester: 7, Age: 22},
			},
			expectedBody: `[{"id":1,"name":"Alice","email":"alice@example.com","department":"CSE","semester":5,"age":21,"createdAtUtc":"0001-01-01T00:00:00Z","updatedAtUtc":"0001-01-01T00:00:00Z"},{"id":2,"name":"Bob","email":"bob@example.com","department":"ECE","semester":3,"age":20,"createdAtUtc":"0001-01-01T00:00:00Z","updatedAtUtc":"0001-01-01T00:00:00Z"},{"id":3,"name":"Charlie","email":"charlie@example.com","department":"IT","semester":7,"age":22,"createdAtUtc":"0001-01-01T00:00:00Z","updatedAtUtc":"0001-01-01T00:00:00Z"}]`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "retrieve all students successful - empty list when no students exist",
			returnEmpty:  true,
			expectedBody: `[]`,
			expectedCode: http.StatusOK,
		},
		{
			name:          "retrieve all students fails due to internal server error",
			mockErr:       mocks.OpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: services.ErrGetAllStudentsFailed.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{
				GetAllStudentsError: tt.mockErr,
				ReturnEmptyStudents: tt.returnEmpty,
				Students:            tt.mockStudents,
			}
			_, mux := setupMockHandler(svc, nil)

			req := httptest.NewRequest(
				http.MethodGet,
				"/students",
				nil,
			)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedCode, rr.Code, rr.Body.String())
			}

			if tt.expectedBody != "" {
				got := rr.Body.String()
				if len(got) > 0 && got[len(got)-1] == '\n' {
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
