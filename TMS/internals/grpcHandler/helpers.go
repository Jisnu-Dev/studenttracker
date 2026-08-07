package grpcHandler

import (
	"net/mail"
	"strings"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxEmailLength = 254

// GrpcError creates and returns a gRPC status error from a status code and error message.
func GrpcError(code codes.Code, message string) error {
	return status.Error(code, message)
}

func GenerateTokenResponse(token string) *tokenpb.GenerateTokenResponse {
	return &tokenpb.GenerateTokenResponse{
		Token: token,
	}
}

func ValidateTokenResponse(isValid bool, adminID int64, adminEmail string) *tokenpb.ValidateTokenResponse {
	return &tokenpb.ValidateTokenResponse{
		IsValid:    isValid,
		AdminId:    adminID,
		AdminEmail: adminEmail,
	}
}

func ValidateGenerateTokenReq(req *tokenpb.GenerateTokenRequest) error {
	if req == nil {
		return GrpcError(codes.InvalidArgument, ErrRequestNil.Error())
	}

	var errs []string

	if req.GetAdminID() <= 0 {
		errs = append(errs, ErrAdminIDInvalid.Error())
	}

	email := strings.TrimSpace(req.GetAdminEmail())
	if email == "" {
		errs = append(errs, ErrAdminEmailRequired.Error())
	} else if len(email) > maxEmailLength {
		errs = append(errs, ErrAdminEmailTooLong.Error())
	} else if strings.ContainsAny(email, " \t\n\r") {
		errs = append(errs, ErrAdminEmailNoWhitespace.Error())
	} else {
		addr, err := mail.ParseAddress(email)
		if err != nil || addr.Address != email || !strings.Contains(email, ".") {
			errs = append(errs, ErrAdminEmailInvalid.Error())
		}
	}

	if len(errs) > 0 {
		return GrpcError(codes.InvalidArgument, strings.Join(errs, "; "))
	}

	return nil
}

func ValidateTokenReq(req *tokenpb.ValidateTokenRequest) error {
	if req == nil {
		return GrpcError(codes.InvalidArgument, ErrRequestNil.Error())
	}

	var errs []string

	token := strings.TrimSpace(req.GetToken())
	if token == "" {
		errs = append(errs, ErrTokenRequired.Error())
	} else {
		if strings.ContainsAny(token, " \t\n\r") {
			errs = append(errs, ErrTokenNoWhitespace.Error())
		}
		if strings.Count(token, ".") != 2 {
			errs = append(errs, ErrTokenMalformed.Error())
		}
	}

	if len(errs) > 0 {
		return GrpcError(codes.InvalidArgument, strings.Join(errs, "; "))
	}

	return nil
}
