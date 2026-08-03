package utils_test

import (
	"testing"

	tokenpb "github.com/Jisnu-Dev/TMS/gen/token"
)

func TestValidateGenerateTokenReq(t *testing.T) {
	tests := []struct {
		name        string
		req         *tokenpb.GenerateTokenRequest
		expectError bool
	}{
		// TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt
		})
	}
}

func TestValidateValidateTokenReq(t *testing.T) {
	tests := []struct {
		name        string
		req         *tokenpb.ValidateTokenRequest
		expectError bool
	}{
		// TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt
		})
	}
}
