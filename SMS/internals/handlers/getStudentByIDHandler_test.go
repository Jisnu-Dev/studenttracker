package handlers_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	mockUtils "github.com/Jisnu-Dev/studenttracker/internals/mocks/utils"
	"github.com/gin-gonic/gin"
)

func TestGetStudentByIDHandler(t *testing.T) {
	tests := []struct {
		name               string
		paramID            string
		mockErr            mocks.MockOpError
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name:               "success - valid id returns student with 200",
			paramID:            "1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"student": map[string]interface{}{
					"id":           float64(1),
					"name":         "test",
					"email":        "test@example.com",
					"department":   "",
					"semester":     float64(0),
					"age":          float64(0),
					"createdAtUtc": "0001-01-01T00:00:00Z",
					"updatedAtUtc": "0001-01-01T00:00:00Z",
				},
			},
		},
		{
			name:               "bad request - non-numeric id param returns 400",
			paramID:            "abc",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "invalid id parameter",
			},
		},
		{
			name:               "not found - student does not exist returns 404",
			paramID:            "99",
			mockErr:            mocks.OpNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"error": "student not found",
			},
		},
		{
			name:               "internal server error - service failure returns 500",
			paramID:            "1",
			mockErr:            mocks.OpInternalError,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "unable to retrieve student",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mockUtils.SetUpGinTestWithParams(
				http.MethodGet, "/students/"+tt.paramID, "",
				gin.Params{{Key: "id", Value: tt.paramID}},
			)

			svc := &mocks.MockService{GetStudentByIDError: tt.mockErr}
			handler := newHandler(svc, nil)
			handler.GetStudentByIDHandler(c)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedStatusCode, w.Code, w.Body.String())
			}

			if tt.expectedResponse != nil {
				var actualResponse map[string]interface{}
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
