package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

func TestGetStudentByIDHandler(t *testing.T) {
	tests := []struct {
		name            string
		paramID         string
		mockErr         mocks.MockOpError
		mockStudent     models.Student
		expectedStudent models.Student
		expectedCode    int
		expectedError   any
	}{
		{
			name:    "retrieve student by id successful",
			paramID: "1",
			mockStudent: models.Student{
				ID:         1,
				Name:       "test",
				Email:      "test@example.com",
				Department: models.CSE,
				Semester:   3,
				Age:        20,
			},
			expectedStudent: models.Student{
				ID:         1,
				Name:       "test",
				Email:      "test@example.com",
				Department: models.CSE,
				Semester:   3,
				Age:        20,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:          "retrieve student fails due to non-numeric id",
			paramID:       "abc",
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidIDParam.Error(),
		},
		{
			name:          "retrieve student fails due to zero id",
			paramID:       "0",
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidIDParam.Error(),
		},
		{
			name:          "retrieve student fails due to negative id",
			paramID:       "-1",
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidIDParam.Error(),
		},
		{
			name:            "retrieve student successful when student not found (returns empty student)",
			paramID:         "99",
			mockErr:         mocks.OpNotFound,
			expectedStudent: models.Student{},
			expectedCode:    http.StatusOK,
		},
		{
			name:          "retrieve student fails due to internal server error",
			paramID:       "1",
			mockErr:       mocks.OpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: services.ErrGetStudentByIDFailed.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{
				GetStudentByIDError: tt.mockErr,
				Student:             tt.mockStudent,
			}
			_, mux := setupMockHandler(svc, nil)

			req := httptest.NewRequest(
				http.MethodGet,
				"/students/"+tt.paramID,
				nil,
			)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedCode, rr.Code, rr.Body.String())
			}

			if tt.expectedCode == http.StatusOK {
				var resp struct {
					Student models.Student `json:"student"`
				}
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if resp.Student.ID != tt.expectedStudent.ID ||
					resp.Student.Name != tt.expectedStudent.Name ||
					resp.Student.Email != tt.expectedStudent.Email ||
					resp.Student.Department != tt.expectedStudent.Department ||
					resp.Student.Semester != tt.expectedStudent.Semester ||
					resp.Student.Age != tt.expectedStudent.Age {
					t.Errorf("expected student %+v, got %+v", tt.expectedStudent, resp.Student)
				}
			}

			if tt.expectedError != nil {
				checkError(t, rr, tt.expectedError)
			}
		})
	}
}
