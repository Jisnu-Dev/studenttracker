package handlers_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	mockUtils "github.com/Jisnu-Dev/studenttracker/internals/mocks/utils"
)

func TestGetAllStudentsHandler(t *testing.T) {
	tests := []struct {
		name               string
		mockErr            mocks.MockOpError
		returnEmpty        bool
		students           []models.Student
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:               "success - returns list of students with 200",
			expectedStatusCode: http.StatusOK,
			expectedResponse: []interface{}{
				map[string]interface{}{
					"id":           float64(1),
					"name":         "test",
					"email":        "test@example.com",
					"department":   "CS",
					"semester":     float64(4),
					"age":          float64(21),
					"createdAtUtc": "0001-01-01T00:00:00Z",
					"updatedAtUtc": "0001-01-01T00:00:00Z",
				},
			},
		},
		{
			name:               "success - returns empty list with 200 when no students exist",
			returnEmpty:        true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   []interface{}{},
		},
		{
			name:               "internal server error - service failure returns 500",
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "unable to retrieve students",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mockUtils.SetUpGinTest(http.MethodGet, "/students", "")

			svc := &mocks.MockService{
				GetAllStudentsError: tt.mockErr,
				ReturnEmptyStudents: tt.returnEmpty,
				Students:            tt.students,
			}
			handler := newHandler(svc, nil)
			handler.GetAllStudentsHandler(c)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedStatusCode, w.Code, w.Body.String())
			}

			if tt.expectedResponse != nil {
				var actualResponse interface{}
				err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
				if err != nil {
					t.Fatalf("failed to parse response body as JSON: %v, body was: %s", err, w.Body.String())
				}

				if !reflect.DeepEqual(tt.expectedResponse, actualResponse) {
					t.Errorf("expected response %v, got %v", tt.expectedResponse, actualResponse)
				}
			}
		})
	}
}
