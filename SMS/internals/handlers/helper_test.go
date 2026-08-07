package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/handlers"
	"github.com/gin-gonic/gin"
)

func TestParseParamID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		paramValue    string
		expectedID    int64
		expectedOk    bool
		expectedCode  int
		expectedError any
	}{
		{
			name:         "valid numeric ID",
			paramValue:   "42",
			expectedID:   42,
			expectedOk:   true,
			expectedCode: http.StatusOK,
		},
		{
			name:          "invalid ID - non numeric",
			paramValue:    "abc",
			expectedOk:    false,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidIDParam.Error(),
		},
		{
			name:          "invalid ID - zero",
			paramValue:    "0",
			expectedOk:    false,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidIDParam.Error(),
		},
		{
			name:          "invalid ID - negative",
			paramValue:    "-10",
			expectedOk:    false,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlers.ErrInvalidIDParam.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tt.paramValue}}

			id, ok := handlers.ParseParamID(c, "id")
			if ok != tt.expectedOk {
				t.Errorf("expected ok %v, got %v", tt.expectedOk, ok)
			}
			if id != tt.expectedID {
				t.Errorf("expected id %d, got %d", tt.expectedID, id)
			}
			if !tt.expectedOk {
				if w.Code != tt.expectedCode {
					t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
				}
				if tt.expectedError != nil {
					checkError(t, w, tt.expectedError)
				}
			}
		})
	}
}
