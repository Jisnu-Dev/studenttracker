package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/grpcClient"
	"github.com/Jisnu-Dev/studenttracker/internals/middlewares"
	"github.com/Jisnu-Dev/studenttracker/internals/mocks"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                 string
		authHeader           string
		mockAdminID          int64
		mockAdminEmail       string
		mockErr              mocks.GrpcOpError
		expectedStatusCode   int
		expectedBody         string
		expectedContextID    int64
		expectedContextEmail string
	}{
		{
			name:                 "success - valid token allows request",
			authHeader:           "Bearer valid.token.here",
			mockAdminID:          1,
			mockAdminEmail:       "admin@example.com",
			expectedStatusCode:   http.StatusOK,
			expectedBody:         `{"message":"success"}`,
			expectedContextID:    1,
			expectedContextEmail: "admin@example.com",
		},
		{
			name:               "unauthorized - missing header",
			authHeader:         "",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"Authorization header is required"}`,
		},
		{
			name:               "unauthorized - missing Bearer prefix",
			authHeader:         "invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"Authorization header format must be 'Bearer \u003ctoken\u003e'"}`,
		},
		{
			name:               "unauthorized - validate token returns false",
			authHeader:         "Bearer invalid.token.here",
			mockErr:            mocks.GrpcOpInvalidToken,
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"Invalid or expired token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			tc := grpcClient.NewTokenClientForTest(&mocks.MockTokenServiceClient{
				ValidateErr: tt.mockErr,
				AdminID:            tt.mockAdminID,
				AdminEmail:         tt.mockAdminEmail,
			})

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

			if got := strings.TrimSpace(w.Body.String()); got != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, got)
			}
		})
	}
}
