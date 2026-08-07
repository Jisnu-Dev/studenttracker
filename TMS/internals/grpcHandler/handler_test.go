package grpcHandler

import (
	"reflect"
	"testing"
)

func TestNewHandler(t *testing.T) {
	secret := "my_super_secret_key"
	handler := NewHandler(secret)

	if handler == nil {
		t.Fatal("expected handler to not be nil")
	}

	expectedJWTSecret := []byte(secret)
	if !reflect.DeepEqual(handler.JWTSecret, expectedJWTSecret) {
		t.Errorf("expected JWTSecret to be %v, got %v", expectedJWTSecret, handler.JWTSecret)
	}
}
