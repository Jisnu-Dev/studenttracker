package clients_test

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	tokenpb "github.com/Jisnu-Dev/studenttracker/gen/token"
	"github.com/Jisnu-Dev/studenttracker/internals/clients"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mockTokenServer is a fake gRPC server used to test the client's behavior.
type mockTokenServer struct {
	tokenpb.UnimplementedTokenServiceServer
	genTokenResp *tokenpb.GenerateTokenResponse
	genTokenErr  error
	valTokenResp *tokenpb.ValidateTokenResponse
	valTokenErr  error
}

func (s *mockTokenServer) GenerateToken(ctx context.Context, req *tokenpb.GenerateTokenRequest) (*tokenpb.GenerateTokenResponse, error) {
	return s.genTokenResp, s.genTokenErr
}

func (s *mockTokenServer) ValidateToken(ctx context.Context, req *tokenpb.ValidateTokenRequest) (*tokenpb.ValidateTokenResponse, error) {
	return s.valTokenResp, s.valTokenErr
}

// setupTestServer starts a mock gRPC server on a random port and returns the address,
// the mock server instance (so tests can modify its behavior), and a cleanup function.
func setupTestServer(t *testing.T) (string, *mockTokenServer, func()) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	mockSrv := &mockTokenServer{}
	tokenpb.RegisterTokenServiceServer(grpcServer, mockSrv)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			// Ignore closed network connection errors during shutdown
			if !strings.Contains(err.Error(), "use of closed network connection") {
				panic("failed to serve test server: " + err.Error())
			}
		}
	}()

	return lis.Addr().String(), mockSrv, func() {
		grpcServer.Stop()
		lis.Close()
	}
}

func TestGenerateTokenClient(t *testing.T) {
	addr, mockSrv, cleanup := setupTestServer(t)
	defer cleanup()

	client, err := clients.NewTokenClient(addr)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		name          string
		adminID       int64
		adminEmail    string
		mockResp      *tokenpb.GenerateTokenResponse
		mockErr       error
		expectedToken string
		expectError   bool
	}{
		{
			name:          "success - returns token on OK",
			adminID:       1,
			adminEmail:    "admin@example.com",
			mockResp:      &tokenpb.GenerateTokenResponse{Token: "mock.jwt.token"},
			mockErr:       nil,
			expectedToken: "mock.jwt.token",
			expectError:   false,
		},
		{
			name:        "error - wraps gRPC server error",
			adminID:     1,
			adminEmail:  "admin@example.com",
			mockResp:    nil,
			mockErr:     status.Error(codes.Internal, "internal server error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSrv.genTokenResp = tt.mockResp
			mockSrv.genTokenErr = tt.mockErr

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			token, err := client.GenerateToken(ctx, tt.adminID, tt.adminEmail)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}

			if err == nil && token != tt.expectedToken {
				t.Errorf("expected token %q, got %q", tt.expectedToken, token)
			}
		})
	}
}

func TestValidateTokenClient(t *testing.T) {
	addr, mockSrv, cleanup := setupTestServer(t)
	defer cleanup()

	client, err := clients.NewTokenClient(addr)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		name               string
		token              string
		mockResp           *tokenpb.ValidateTokenResponse
		mockErr            error
		expectedValid      bool
		expectedAdminID    int64
		expectedAdminEmail string
		expectError        bool
	}{
		{
			name:  "success - valid token returns true and admin details",
			token: "mock.valid.token",
			mockResp: &tokenpb.ValidateTokenResponse{
				IsValid:    true,
				AdminID:    1,
				AdminEmail: "admin@example.com",
			},
			mockErr:            nil,
			expectedValid:      true,
			expectedAdminID:    1,
			expectedAdminEmail: "admin@example.com",
			expectError:        false,
		},
		{
			name:  "success - invalid token returns false and zero values",
			token: "mock.invalid.token",
			mockResp: &tokenpb.ValidateTokenResponse{
				IsValid:    false,
				AdminID:    0,
				AdminEmail: "",
			},
			mockErr:            nil,
			expectedValid:      false,
			expectedAdminID:    0,
			expectedAdminEmail: "",
			expectError:        false,
		},
		{
			name:        "error - wraps gRPC server error",
			token:       "mock.error.token",
			mockResp:    nil,
			mockErr:     status.Error(codes.Internal, "internal server error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSrv.valTokenResp = tt.mockResp
			mockSrv.valTokenErr = tt.mockErr

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			isValid, adminID, adminEmail, err := client.ValidateToken(ctx, tt.token)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}

			if err == nil {
				if isValid != tt.expectedValid {
					t.Errorf("expected isValid=%v, got %v", tt.expectedValid, isValid)
				}
				if tt.expectedValid {
					if adminID != tt.expectedAdminID {
						t.Errorf("expected adminID=%d, got %d", tt.expectedAdminID, adminID)
					}
					if adminEmail != tt.expectedAdminEmail {
						t.Errorf("expected adminEmail=%q, got %q", tt.expectedAdminEmail, adminEmail)
					}
				}
			}
		})
	}
}
