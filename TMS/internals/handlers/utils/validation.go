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
	if req.GetAdminID() <= 0 {
		return GrpcError(codes.InvalidArgument, "admin_id is required and must be greater than 0")
	}

	email := strings.TrimSpace(req.GetAdminEmail())
	if email == "" {
		return GrpcError(codes.InvalidArgument, "admin_email is required")
	}
	if len(email) > maxEmailLength {
		return GrpcError(codes.InvalidArgument, "admin_email exceeds maximum length")
	}
	if strings.ContainsAny(email, " \t\n\r") {
		return GrpcError(codes.InvalidArgument, "admin_email cannot contain whitespace")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email || !strings.Contains(email, ".") {
		return GrpcError(codes.InvalidArgument, "admin_email is invalid")
	}

	return nil
}

func ValidateValidateTokenReq(req *tokenpb.ValidateTokenRequest) error {
	if req == nil {
		return GrpcError(codes.InvalidArgument, "request cannot be nil")
	}

	token := strings.TrimSpace(req.GetToken())
	if token == "" {
		return GrpcError(codes.InvalidArgument, "token is required")
	}
	if strings.ContainsAny(token, " \t\n\r") {
		return GrpcError(codes.InvalidArgument, "token cannot contain whitespace")
	}

	if strings.Count(token, ".") != 2 {
		return GrpcError(codes.InvalidArgument, "malformed token structure")
	}

	return nil
}
