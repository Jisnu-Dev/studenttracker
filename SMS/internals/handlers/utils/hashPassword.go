package utils

import (
	"golang.org/x/crypto/bcrypt"
)

var HashPassword = func(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
