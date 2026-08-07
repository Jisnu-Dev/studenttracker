package handlers_test

import (
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
		expectedBody  string
		expectedError any
	}{
		{
			name:         "delete student successful",
			paramID:      "1",
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"message":"Student deleted successfully"}`,
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
			expectedBody: `{"id":99,"message":"Student deleted successfully"}`,
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
