package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

func TestDeleteStudentHandler(t *testing.T) {
	tests := []struct {
		name          string
		paramID       string
		mockErr       mocks.MockOpError
		expectedCode  int
		expectedError any
	}{
		{
			name:         "delete student successful",
			paramID:      "1",
			expectedCode: http.StatusOK,
		},
		{
			name:          "delete student fails due to non-numeric id",
			paramID:       "abc",
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidIDParam.Error(),
		},
		{
			name:          "delete student fails due to float id",
			paramID:       "1.5",
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidIDParam.Error(),
		},
		{
			name:          "delete student fails due to zero id",
			paramID:       "0",
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidIDParam.Error(),
		},
		{
			name:          "delete student fails due to negative id",
			paramID:       "-1",
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidIDParam.Error(),
		},
		{
			name:         "delete student successful when student does not exist",
			paramID:      "99",
			mockErr:      mocks.OpNotFound,
			expectedCode: http.StatusOK,
		},
		{
			name:          "delete student fails due to internal server error",
			paramID:       "1",
			mockErr:       mocks.OpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: services.ErrDeleteStudentFailed.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{DeleteStudentError: tt.mockErr}
			_, mux := setupMockHandler(svc, nil)

			req := httptest.NewRequest(
				http.MethodDelete,
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
					ID      int64  `json:"id"`
					Message string `json:"message"`
				}
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if resp.Message != "Student deleted successfully" {
					t.Errorf("expected message 'Student deleted successfully', got %q", resp.Message)
				}
			}

			if tt.expectedError != nil {
				checkError(t, rr, tt.expectedError)
			}
		})
	}
}
