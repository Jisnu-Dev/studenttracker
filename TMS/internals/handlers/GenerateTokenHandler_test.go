package handlers_test

import (
	"context"
	"testing"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/Jisnu-Dev/TMS/internals/handlers"
	"github.com/Jisnu-Dev/TMS/internals/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newTMSHandler(svc *mocks.MockService) *handlers.Handler {
	return handlers.NewHandler(svc)
}

func TestGenerateTokenHandler(t *testing.T) {
	tests := []struct {
		name             string
		request          *tokenpb.GenerateTokenRequest
		mockErr          mocks.MockOpError
		expectedCode     codes.Code
		expectedResponse *tokenpb.GenerateTokenResponse
	}{
		{
			name: "success - valid request returns token with OK",
			request: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "admin@example.com",
			},
			expectedCode: codes.OK,
			expectedResponse: &tokenpb.GenerateTokenResponse{
				Token: "mocked.jwt.token",
			},
		},
		{
			name:         "invalid argument - nil request returns InvalidArgument",
			request:      nil,
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "invalid argument - admin_id is zero returns InvalidArgument",
			request: &tokenpb.GenerateTokenRequest{
				AdminID:    0,
				AdminEmail: "admin@example.com",
			},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "invalid argument - negative admin_id returns InvalidArgument",
			request: &tokenpb.GenerateTokenRequest{
				AdminID:    -1,
				AdminEmail: "admin@example.com",
			},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "invalid argument - empty admin_email returns InvalidArgument",
			request: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "",
			},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "invalid argument - invalid admin_email format returns InvalidArgument",
			request: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "not-an-email",
			},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "internal - service failure returns Internal",
			request: &tokenpb.GenerateTokenRequest{
				AdminID:    1,
				AdminEmail: "admin@example.com",
			},
			mockErr:      mocks.OpInternalError,
			expectedCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{GenerateTokenError: tt.mockErr}
			handler := newTMSHandler(svc)

			resp, err := handler.GenerateToken(context.Background(), tt.request)

			if tt.expectedCode == codes.OK {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if resp == nil {
					t.Fatal("expected non-nil response")
				}

				// Compare pb response values using reflect.DeepEqual (but ensure no inner proto state differs)
				// For simple token string response, direct check is also fine, but checking expectedResponse fields explicitly.
				if tt.expectedResponse != nil {
					if resp.GetToken() != tt.expectedResponse.GetToken() {
						t.Errorf("expected response %v, got %v", tt.expectedResponse, resp)
					}
				}
				return
			}

			if err == nil {
				t.Fatal("expected an error but got nil")
			}
			st, ok := status.FromError(err)
			if !ok {
				t.Fatalf("expected a gRPC status error, got %v", err)
			}
			if st.Code() != tt.expectedCode {
				t.Errorf("expected gRPC code %v, got %v (message: %s)", tt.expectedCode, st.Code(), st.Message())
			}
		})
	}
}
