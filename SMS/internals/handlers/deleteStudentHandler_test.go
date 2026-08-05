package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
)

func TestDeleteStudentHandler(t *testing.T) {
	tests := []struct {
		name               string
		paramID            string
		mockErr            mocks.MockOpError
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "success - valid id deletes student and returns 200",
			paramID:            "1",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"id":1,"message":"Student deleted successfully"}`,
		},
		{
			name:               "bad request - non-numeric id param returns 400",
			paramID:            "abc",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid id parameter"}`,
		},
		{
			name:               "bad request - float id param returns 400",
			paramID:            "1.5",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid id parameter"}`,
		},
		{
			name:               "bad request - zero id param returns 400",
			paramID:            "0",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid id parameter"}`,
		},
		{
			name:               "bad request - negative id param returns 400",
			paramID:            "-1",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"invalid id parameter"}`,
		},
		{
			name:               "success - valid non-existing id returns 200",
			paramID:            "99",
			mockErr:            mocks.OpNotFound,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"id":99,"message":"Student deleted successfully"}`,
		},
		{
			name:               "internal server error - service failure returns 500",
			paramID:            "1",
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"unable to delete student"}`,
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

			if rr.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedStatusCode, rr.Code, rr.Body.String())
			}

			if got := strings.TrimSpace(rr.Body.String()); got != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, got)
			}
		})
	}
}
