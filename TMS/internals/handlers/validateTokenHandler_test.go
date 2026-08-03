package handlers_test

import (
	"context"
	"testing"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"github.com/Jisnu-Dev/TMS/internals/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidateTokenHandler(t *testing.T) {
	tests := []struct {
		name             string
		request          *tokenpb.ValidateTokenRequest
		mockErr          mocks.MockOpError
		expectedCode     codes.Code
		expectedResponse *tokenpb.ValidateTokenResponse
	}{
		{
			name:             "success - valid token returns isValid=true with admin details",
			request:          &tokenpb.ValidateTokenRequest{Token: "mock.valid.token"},
			expectedCode:     codes.OK,
			expectedResponse: &tokenpb.ValidateTokenResponse{
				IsValid:    true,
				AdminId:    1,
				AdminEmail: "admin@example.com",
			},
		},
		{
			name:          "invalid argument - nil request returns InvalidArgument",
			request:       nil,
			expectedCode:  codes.InvalidArgument,
		},
		{
			name:          "invalid argument - empty token returns InvalidArgument",
			request:       &tokenpb.ValidateTokenRequest{Token: ""},
			expectedCode:  codes.InvalidArgument,
		},
		{
			name:          "invalid argument - token with whitespace returns InvalidArgument",
			request:       &tokenpb.ValidateTokenRequest{Token: "mock valid token"},
			expectedCode:  codes.InvalidArgument,
		},
		{
			name:          "invalid argument - token with wrong number of segments returns InvalidArgument",
			request:       &tokenpb.ValidateTokenRequest{Token: "onlyone"},
			expectedCode:  codes.InvalidArgument,
		},
		{
			name:             "success - invalid token from service returns isValid=false with zero values",
			request:          &tokenpb.ValidateTokenRequest{Token: "mock.invalid.token"},
			mockErr:          mocks.OpInvalidToken,
			expectedCode:     codes.OK,
			expectedResponse: &tokenpb.ValidateTokenResponse{
				IsValid:    false,
				AdminId:    0,
				AdminEmail: "",
			},
		},
		{
			name:             "success - signing method mismatch returns isValid=false with zero values",
			request:          &tokenpb.ValidateTokenRequest{Token: "mock.mismatch.token"},
			mockErr:          mocks.OpSigningMethodMismatch,
			expectedCode:     codes.OK,
			expectedResponse: &tokenpb.ValidateTokenResponse{
				IsValid:    false,
				AdminId:    0,
				AdminEmail: "",
			},
		},
		{
			name:             "success - internal service error returns isValid=false with zero values",
			request:          &tokenpb.ValidateTokenRequest{Token: "mock.error.token"},
			mockErr:          mocks.OpInternalError,
			expectedCode:     codes.OK,
			expectedResponse: &tokenpb.ValidateTokenResponse{
				IsValid:    false,
				AdminId:    0,
				AdminEmail: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockService{ValidateTokenError: tt.mockErr}
			handler := newTMSHandler(svc)

			resp, err := handler.ValidateToken(context.Background(), tt.request)

			if tt.expectedCode == codes.InvalidArgument {
				if err == nil {
					t.Fatal("expected an error but got nil")
				}
				st, ok := status.FromError(err)
				if !ok {
					t.Fatalf("expected a gRPC status error, got %v", err)
				}
				if st.Code() != codes.InvalidArgument {
					t.Errorf("expected gRPC code %v, got %v", codes.InvalidArgument, st.Code())
				}
				return
			}

			// For codes.OK cases (which includes all service-level failures — they return isValid=false, not a gRPC error)
			if err != nil {
				t.Fatalf("expected no gRPC error, got %v", err)
			}
			if resp == nil {
				t.Fatal("expected non-nil response")
			}
			
			if tt.expectedResponse != nil {
				if resp.GetIsValid() != tt.expectedResponse.GetIsValid() {
					t.Errorf("expected isValid=%v, got %v", tt.expectedResponse.GetIsValid(), resp.GetIsValid())
				}
				if resp.GetAdminId() != tt.expectedResponse.GetAdminId() {
					t.Errorf("expected admin_id=%d, got %d", tt.expectedResponse.GetAdminId(), resp.GetAdminId())
				}
				if resp.GetAdminEmail() != tt.expectedResponse.GetAdminEmail() {
					t.Errorf("expected admin_email=%q, got %q", tt.expectedResponse.GetAdminEmail(), resp.GetAdminEmail())
				}
			}
		})
	}
}
