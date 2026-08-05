package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
)

func TestGetStudentByIDHandler(t *testing.T) {
	tests := []struct {
		name               string
		paramID            string
		mockErr            mocks.MockOpError
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "success - valid id returns student with 200",
			paramID:            "1",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"student":{"id":1,"name":"test","email":"test@example.com","department":"","semester":0,"age":0,"createdAtUtc":"0001-01-01T00:00:00Z","updatedAtUtc":"0001-01-01T00:00:00Z"}}`,
		},
		{
			name:               "bad request - non-numeric id param returns 400",
			paramID:            "abc",
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
			name:               "success - valid id not present returns empty object with 200",
			paramID:            "99",
			mockErr:            mocks.OpNotFound,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{}`,
		},
		{
			name:               "internal server error - service failure returns 500",
			paramID:            "1",
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"unable to retrieve student"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{GetStudentByIDError: tt.mockErr}
			_, mux := setupMockHandler(svc, nil)

			req := httptest.NewRequest(
				http.MethodGet,
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
