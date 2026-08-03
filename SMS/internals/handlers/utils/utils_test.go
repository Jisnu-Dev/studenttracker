package utils_test

import (
	"testing"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		checkPassword string
		expectMatch   bool
		expectError   bool
	}{
		// TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt
		})
	}
}
