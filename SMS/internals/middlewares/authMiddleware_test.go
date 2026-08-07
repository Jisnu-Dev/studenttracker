package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		returnInvalid        bool
		expectedCode         int
		expectedError        string
		expectedContextID    int64
		expectedContextEmail string
	}{
		{
			name:                 "valid token allows request",
			authHeader:           "Bearer valid.token.here",
			mockAdminID:          1,
			mockAdminEmail:       "admin@example.com",
			expectedCode:         http.StatusOK,
			expectedContextID:    1,
			expectedContextEmail: "admin@example.com",
		},
		{
			name:          "unauthorized - missing header",
			authHeader:    "",
			expectedCode:  http.StatusUnauthorized,
			expectedError: middlewares.ErrAuthHeaderRequired.Error(),
		},
		{
			name:          "unauthorized - missing Bearer prefix",
			authHeader:    "invalid.token.here",
			expectedCode:  http.StatusUnauthorized,
			expectedError: middlewares.ErrAuthHeaderFormat.Error(),
		},
		{
			name:          "unauthorized - validate token returns gRPC error",
			authHeader:    "Bearer invalid.token.here",
			mockErr:       mocks.GrpcOpInvalidToken,
			expectedCode:  http.StatusUnauthorized,
			expectedError: middlewares.ErrInvalidToken.Error(),
		},
		{
			name:          "unauthorized - validate token reports token invalid",
			authHeader:    "Bearer expired.token.here",
			returnInvalid: true,
			expectedCode:  http.StatusUnauthorized,
			expectedError: middlewares.ErrInvalidToken.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			tc := grpcClient.NewTokenClientForTest(&mocks.MockTokenServiceClient{
				ValidateTokenError: tt.mockErr,
				ReturnInvalidToken: tt.returnInvalid,
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

			if w.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d (body: %s)", tt.expectedCode, w.Code, w.Body.String())
			}

			if tt.expectedCode == http.StatusOK {
				var resp struct {
					Message string `json:"message"`
				}
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if resp.Message != "success" {
					t.Errorf("expected message 'success', got %q", resp.Message)
				}
			}

			if tt.expectedError != "" {
				var resp struct {
					Error string `json:"error"`
				}
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode error body: %v", err)
				}
				if resp.Error != tt.expectedError {
					t.Errorf("expected error %q, got %q", tt.expectedError, resp.Error)
				}
			}
		})
	}
}
