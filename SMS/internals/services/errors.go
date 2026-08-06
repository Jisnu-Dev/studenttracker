package services

import "errors"

var (
	// Student Domain Errors
	ErrStudentNotFound    = errors.New("student not found")
	ErrStudentEmailExists = errors.New("student with this email already exists")

	// Admin & Auth Domain Errors
	ErrAdminNotFound      = errors.New("admin not found")
	ErrAdminEmailExists   = errors.New("admin with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")

	// Request Errors
	ErrNoFieldsToUpdate = errors.New("at least one field must be provided for update")

	// Fallback Function Execution Errors
	ErrCreateStudentFailed   = errors.New("unable to create student")
	ErrGetAllStudentsFailed  = errors.New("unable to retrieve students")
	ErrGetStudentByIDFailed  = errors.New("unable to retrieve student")
	ErrDeleteStudentFailed   = errors.New("unable to delete student")
	ErrUpdateStudentFailed   = errors.New("unable to update student")
	ErrPatchStudentFailed    = errors.New("unable to patch student")
	ErrRegisterAdminFailed   = errors.New("unable to register admin")
	ErrGetAdminByEmailFailed = errors.New("unable to retrieve admin")
)
