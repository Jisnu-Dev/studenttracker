package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

func TestGetAllStudentsHandler(t *testing.T) {
	tests := []struct {
		name             string
		mockErr          mocks.MockOpError
		returnEmpty      bool
		mockStudents     []models.Student
		expectedStudents []models.Student
		expectedCode     int
		expectedError    any
	}{
		{
			name: "retrieve all students successful - single student",
			mockStudents: []models.Student{
				{ID: 1, Name: "test", Email: "test@example.com", Department: "CS", Semester: 4, Age: 21},
			},
			expectedStudents: []models.Student{
				{ID: 1, Name: "test", Email: "test@example.com", Department: "CS", Semester: 4, Age: 21},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "retrieve all students successful - multiple students",
			mockStudents: []models.Student{
				{ID: 1, Name: "Alice", Email: "alice@example.com", Department: models.CSE, Semester: 5, Age: 21},
				{ID: 2, Name: "Bob", Email: "bob@example.com", Department: models.ECE, Semester: 3, Age: 20},
				{ID: 3, Name: "Charlie", Email: "charlie@example.com", Department: models.IT, Semester: 7, Age: 22},
			},
			expectedStudents: []models.Student{
				{ID: 1, Name: "Alice", Email: "alice@example.com", Department: models.CSE, Semester: 5, Age: 21},
				{ID: 2, Name: "Bob", Email: "bob@example.com", Department: models.ECE, Semester: 3, Age: 20},
				{ID: 3, Name: "Charlie", Email: "charlie@example.com", Department: models.IT, Semester: 7, Age: 22},
			},
			expectedCode: http.StatusOK,
		},
		{
			name:             "retrieve all students successful - empty list when no students exist",
			returnEmpty:      true,
			expectedStudents: []models.Student{},
			expectedCode:     http.StatusOK,
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

			if tt.expectedCode == http.StatusOK {
				var gotStudents []models.Student
				if err := json.NewDecoder(rr.Body).Decode(&gotStudents); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if len(gotStudents) != len(tt.expectedStudents) {
					t.Fatalf("expected %d students, got %d", len(tt.expectedStudents), len(gotStudents))
				}
				for i := range gotStudents {
					if gotStudents[i].ID != tt.expectedStudents[i].ID ||
						gotStudents[i].Name != tt.expectedStudents[i].Name ||
						gotStudents[i].Email != tt.expectedStudents[i].Email ||
						gotStudents[i].Department != tt.expectedStudents[i].Department ||
						gotStudents[i].Semester != tt.expectedStudents[i].Semester ||
						gotStudents[i].Age != tt.expectedStudents[i].Age {
						t.Errorf("student at index %d mismatch: expected %+v, got %+v", i, tt.expectedStudents[i], gotStudents[i])
					}
				}
			}

			if tt.expectedError != nil {
				checkError(t, rr, tt.expectedError)
			}
		})
	}
}
