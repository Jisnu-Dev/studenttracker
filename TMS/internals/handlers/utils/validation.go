package utils

import (
	"net/mail"
	"strings"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
	"google.golang.org/grpc/codes"
)

const maxEmailLength = 254

func ValidateGenerateTokenReq(req *tokenpb.GenerateTokenRequest) error {
	if req == nil {
		return GrpcError(codes.InvalidArgument, "request cannot be nil")
	}

	var errs []string

	if req.GetAdminID() <= 0 {
		errs = append(errs, "admin_id is required and must be greater than 0")
	}

	email := strings.TrimSpace(req.GetAdminEmail())
	if email == "" {
		errs = append(errs, "admin_email is required")
	} else if len(email) > maxEmailLength {
		errs = append(errs, "admin_email exceeds maximum length")
	} else if strings.ContainsAny(email, " \t\n\r") {
		errs = append(errs, "admin_email cannot contain whitespace")
	} else {
		addr, err := mail.ParseAddress(email)
		if err != nil || addr.Address != email || !strings.Contains(email, ".") {
			errs = append(errs, "admin_email is invalid")
		}
	}

	if len(errs) > 0 {
		return GrpcError(codes.InvalidArgument, strings.Join(errs, "; "))
	}

	return nil
}

func ValidateValidateTokenReq(req *tokenpb.ValidateTokenRequest) error {
	if req == nil {
		return GrpcError(codes.InvalidArgument, "request cannot be nil")
	}

	var errs []string

	token := strings.TrimSpace(req.GetToken())
	if token == "" {
		errs = append(errs, "token is required")
	} else {
		if strings.ContainsAny(token, " \t\n\r") {
			errs = append(errs, "token cannot contain whitespace")
		}
		if strings.Count(token, ".") != 2 {
			errs = append(errs, "malformed token structure")
		}
	}

	if len(errs) > 0 {
		return GrpcError(codes.InvalidArgument, strings.Join(errs, "; "))
	}

	return nil
}
