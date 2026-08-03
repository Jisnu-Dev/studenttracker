package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/middlewares"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		authHeader         string
		mockAdminID        int64
		mockAdminEmail     string
		mockErr            mocks.MockOpError
		expectedStatusCode int
		expectedResponse   map[string]interface{}
		expectedContextID  int64
		expectedContextEmail string
	}{
		{
			name:               "success - valid token allows request",
			authHeader:         "Bearer valid.token.here",
			mockAdminID:        1,
			mockAdminEmail:     "admin@example.com",
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"message": "success",
			},
			expectedContextID:    1,
			expectedContextEmail: "admin@example.com",
		},
		{
			name:               "unauthorized - missing header",
			authHeader:         "",
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: map[string]interface{}{
				"error": "Authorization header is required",
			},
		},
		{
			name:               "unauthorized - missing Bearer prefix",
			authHeader:         "invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: map[string]interface{}{
				"error": "Authorization header format must be 'Bearer <token>'",
			},
		},
		{
			name:               "unauthorized - validate token returns false",
			authHeader:         "Bearer invalid.token.here",
			mockErr:            mocks.OpInvalidToken,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: map[string]interface{}{
				"error": "Invalid or expired token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			tc := &mocks.MockTokenClient{
				ValidateTokenError: tt.mockErr,
				AdminID:            tt.mockAdminID,
				Email:              tt.mockAdminEmail,
			}

			// Add a dummy route to test if middleware allows request to pass
			r.GET("/protected", middlewares.AuthMiddleware(tc), func(ctx *gin.Context) {
				adminID, _ := ctx.Get("adminID")
				adminEmail, _ := ctx.Get("adminEmail")
				
				if adminID != tt.expectedContextID {
					t.Errorf("expected context adminID %v, got %v", tt.expectedContextID, adminID)
				}
				if adminEmail != tt.expectedContextEmail {
					t.Errorf("expected context adminEmail %v, got %v", tt.expectedContextEmail, adminEmail)
				}

				ctx.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			c.Request = req

			r.ServeHTTP(w, req)

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
