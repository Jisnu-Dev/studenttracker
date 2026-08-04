package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
)

func TestGetAllStudentsHandler(t *testing.T) {
	tests := []struct {
		name               string
		mockErr            mocks.MockOpError
		returnEmpty        bool
		students           []models.Student
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "success - returns list of students with 200",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `[{"id":1,"name":"test","email":"test@example.com","department":"CS","semester":4,"age":21,"createdAtUtc":"0001-01-01T00:00:00Z","updatedAtUtc":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name:               "success - returns empty list with 200 when no students exist",
			returnEmpty:        true,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `[]`,
		},
		{
			name:               "internal server error - service failure returns 500",
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"unable to retrieve students"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{
				GetAllStudentsError: tt.mockErr,
				ReturnEmptyStudents: tt.returnEmpty,
				Students:            tt.students,
			}
			_, mux := setupMockHandler(svc, nil)

			req := httptest.NewRequest(
				http.MethodGet,
				"/students",
				nil,
			)
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
