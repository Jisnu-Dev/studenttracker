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
			name:               "success - returns single student with 200",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `[{"id":1,"name":"test","email":"test@example.com","department":"CS","semester":4,"age":21,"createdAtUtc":"0001-01-01T00:00:00Z","updatedAtUtc":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "success - returns multiple students with 200",
			students: []models.Student{
				{
					ID:         1,
					Name:       "Alice",
					Email:      "alice@example.com",
					Department: models.CSE,
					Semester:   5,
					Age:        21,
				},
				{
					ID:         2,
					Name:       "Bob",
					Email:      "bob@example.com",
					Department: models.ECE,
					Semester:   3,
					Age:        20,
				},
				{
					ID:         3,
					Name:       "Charlie",
					Email:      "charlie@example.com",
					Department: models.IT,
					Semester:   7,
					Age:        22,
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `[{"id":1,"name":"Alice","email":"alice@example.com","department":"CSE","semester":5,"age":21,"createdAtUtc":"0001-01-01T00:00:00Z","updatedAtUtc":"0001-01-01T00:00:00Z"},{"id":2,"name":"Bob","email":"bob@example.com","department":"ECE","semester":3,"age":20,"createdAtUtc":"0001-01-01T00:00:00Z","updatedAtUtc":"0001-01-01T00:00:00Z"},{"id":3,"name":"Charlie","email":"charlie@example.com","department":"IT","semester":7,"age":22,"createdAtUtc":"0001-01-01T00:00:00Z","updatedAtUtc":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name:               "success - returns empty object with 200 when no students exist",
			returnEmpty:        true,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{}`,
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
