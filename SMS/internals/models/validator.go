package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

var validDepartments = map[DepartmentType]bool{
	CSE:   true,
	IT:    true,
	ECE:   true,
	EEE:   true,
	MECH:  true,
	CIVIL: true,
	AIDS:  true,
	AIML:  true,
}

func init() {
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	Validate.RegisterValidation("name", validateName)
	Validate.RegisterValidation("department", validateDepartment)
	Validate.RegisterValidation("password", validatePassword)
}

func validateDepartment(fl validator.FieldLevel) bool {
	dept := DepartmentType(fl.Field().String())
	return validDepartments[dept]
}

func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	if strings.TrimSpace(name) != name {
		return false
	}

	if strings.Contains(name, "  ") {
		return false
	}

	for _, r := range name {
		switch {
		case unicode.IsLetter(r):
			continue

		case r == ' ':
			continue

		default:
			return false
		}
	}

	return true
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)
	const specialChars = "!@#$%^&*()-_=+[]{}|\\:;\"'<>,.?/~`"
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true

		case unicode.IsLower(r):
			hasLower = true

		case unicode.IsDigit(r):
			hasDigit = true

		case unicode.IsSpace(r):
			return false

		case strings.ContainsRune(specialChars, r):
			hasSpecial = true
		}
	}
	return hasUpper &&
		hasLower &&
		hasDigit &&
		hasSpecial
}

func fieldError(e validator.FieldError) error {
	switch e.Field() {

	case "name":
		switch e.Tag() {
		case "required":
			return ErrNameRequired
		case "min":
			return ErrNameTooShort
		case "max":
			return ErrNameTooLong
		default:
			return ErrInvalidName
		}

	case "email":
		switch e.Tag() {
		case "required":
			return ErrEmailRequired
		case "max":
			return ErrEmailTooLong
		default:
			return ErrInvalidEmail
		}

	case "department":
		switch e.Tag() {
		case "required":
			return ErrDepartmentRequired
		default:
			return ErrInvalidDepartment
		}

	case "semester":
		return ErrSemesterInvalid

	case "age":
		return ErrAgeInvalid

	case "password":
		switch e.Tag() {
		case "required":
			return ErrPasswordRequired
		case "min":
			return ErrPasswordTooShort
		case "max":
			return ErrPasswordTooLong
		default:
			return ErrInvalidPassword
		}

	default:
		return fmt.Errorf("%s failed validation on '%s'", e.Field(), e.Tag())
	}
}

func ValidationErrors(err error) map[string]string {
	var validationErrors validator.ValidationErrors

	if !errors.As(err, &validationErrors) {
		return map[string]string{
			"error": err.Error(),
		}
	}

	result := make(map[string]string)

	for _, e := range validationErrors {
		result[e.Field()] = fieldError(e).Error()
	}

	return result
}
