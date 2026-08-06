package models

import "errors"

var (
	ErrNameRequired = errors.New("name is required")
	ErrNameTooShort = errors.New("name must be at least 2 characters")
	ErrNameTooLong  = errors.New("name must be at most 100 characters")
	ErrInvalidName  = errors.New("name may only contain letters separated by single spaces")

	ErrEmailRequired = errors.New("email is required")
	ErrEmailTooLong  = errors.New("email cannot exceed 254 characters")
	ErrInvalidEmail  = errors.New("email is invalid")

	ErrDepartmentRequired = errors.New("department is required")
	ErrInvalidDepartment  = errors.New("invalid department")

	ErrSemesterInvalid = errors.New("semester must be between 1 and 8")

	ErrAgeInvalid = errors.New("age must be between 18 and 60")

	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong  = errors.New("password cannot exceed 64 characters")
	ErrInvalidPassword  = errors.New("password must contain uppercase, lowercase, digit, and special character")

	ErrAtLeastOneField = errors.New("at least one field must be provided for update")
)
