package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

func TestGetStudentByIDHandler(t *testing.T) {
	tests := []struct {
		name          string
		paramID       string
		mockErr       mocks.MockOpError
		expectedBody  string
		expectedCode  int
		expectedError any
	}{
		{
			name:         "retrieve student by id successful",
			paramID:      "1",
			expectedBody: `{"student":{"id":1,"name":"test","email":"test@example.com","department":"CSE","semester":3,"age":20,"createdAtUtc":"0001-01-01T00:00:00Z","updatedAtUtc":"0001-01-01T00:00:00Z"}}`,
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
			name:         "retrieve student successful when student not found (returns empty student)",
			paramID:      "99",
			mockErr:      mocks.OpNotFound,
			expectedBody: `[]`,
			expectedCode: http.StatusOK,
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
