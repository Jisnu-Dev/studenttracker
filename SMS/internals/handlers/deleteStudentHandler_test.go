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

func TestDeleteStudentHandler(t *testing.T) {
	tests := []struct {
		name               string
		paramID            string
		mockErr            mocks.MockOpError
		expectedStatusCode int
		expectedResponse   map[string]interface{}
	}{
		{
			name:               "success - valid id deletes student and returns 200",
			paramID:            "1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"id":      float64(1), // JSON numbers parse as float64
				"message": "Student deleted successfully",
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
			name:               "bad request - float id param returns 400",
			paramID:            "1.5",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "invalid id parameter",
			},
		},
		{
			name:               "not found - zero id returns 404 (does not exist)",
			paramID:            "0",
			mockErr:            mocks.OpNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"error": "student not found",
			},
		},
		{
			name:               "not found - negative id returns 404 (does not exist)",
			paramID:            "-1",
			mockErr:            mocks.OpNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"error": "student not found",
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
				"error": "unable to delete student",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mockUtils.SetUpGinTestWithParams(
				http.MethodDelete, "/students/"+tt.paramID, "",
				gin.Params{{Key: "id", Value: tt.paramID}},
			)

			svc := &mocks.MockService{DeleteStudentError: tt.mockErr}
			handler := newHandler(svc, nil)
			handler.DeleteStudentHandler(c)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedStatusCode, w.Code, w.Body.String())
			}

			if tt.expectedResponse != nil {
				var actualResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
				if err != nil {
					t.Fatalf("failed to parse response body as JSON: %v, body was: %s", err, w.Body.String())
				}

				for key, expectedValue := range tt.expectedResponse {
					actualValue, exists := actualResponse[key]
					if !exists {
						t.Errorf("expected key %q in response, but it was missing. Full response: %v", key, actualResponse)
						continue
					}
					if !reflect.DeepEqual(expectedValue, actualValue) {
						t.Errorf("expected %q to be %v, got %v", key, expectedValue, actualValue)
					}
				}
			}
		})
	}
}
