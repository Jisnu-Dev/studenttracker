package models

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateStudent(t *testing.T) {
	validStudent := func() Student {
		return Student{
			Name:       "test name",
			Email:      "test@example.com",
			Department: CSE,
			Semester:   4,
			Age:        20,
		}
	}

	t.Run("valid student passes", func(t *testing.T) {
		s := validStudent()
		if err := Validate.Struct(s); err != nil {
			t.Fatalf("expected valid student to pass, got error: %v", err)
		}
	})

	t.Run("missing name returns required error", func(t *testing.T) {
		s := validStudent()
		s.Name = ""
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["name"] != ErrNameRequired.Error() {
			t.Errorf("expected %q, got %q", ErrNameRequired.Error(), errs["name"])
		}
	})

	t.Run("name with invalid characters returns invalid name error", func(t *testing.T) {
		s := validStudent()
		s.Name = "Test123"
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["name"] != ErrInvalidName.Error() {
			t.Errorf("expected %q, got %q", ErrInvalidName.Error(), errs["name"])
		}
	})

	t.Run("name too short returns min error", func(t *testing.T) {
		s := validStudent()
		s.Name = "J"
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["name"] != ErrNameTooShort.Error() {
			t.Errorf("expected %q, got %q", ErrNameTooShort.Error(), errs["name"])
		}
	})

	t.Run("name too long returns max error", func(t *testing.T) {
		s := validStudent()
		s.Name = strings.Repeat("a", 101)
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["name"] != ErrNameTooLong.Error() {
			t.Errorf("expected %q, got %q", ErrNameTooLong.Error(), errs["name"])
		}
	})

	t.Run("name with leading or trailing spaces returns invalid name error", func(t *testing.T) {
		s := validStudent()
		s.Name = " Test Name "
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["name"] != ErrInvalidName.Error() {
			t.Errorf("expected %q, got %q", ErrInvalidName.Error(), errs["name"])
		}
	})

	t.Run("name with consecutive spaces returns invalid name error", func(t *testing.T) {
		s := validStudent()
		s.Name = "Test  Name"
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["name"] != ErrInvalidName.Error() {
			t.Errorf("expected %q, got %q", ErrInvalidName.Error(), errs["name"])
		}
	})

	t.Run("missing email returns required error", func(t *testing.T) {
		s := validStudent()
		s.Email = ""
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["email"] != ErrEmailRequired.Error() {
			t.Errorf("expected %q, got %q", ErrEmailRequired.Error(), errs["email"])
		}
	})

	t.Run("email too long returns max error", func(t *testing.T) {
		s := validStudent()
		s.Email = strings.Repeat("a", 250) + "@example.com"
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["email"] != ErrEmailTooLong.Error() {
			t.Errorf("expected %q, got %q", ErrEmailTooLong.Error(), errs["email"])
		}
	})

	t.Run("invalid email returns invalid email error", func(t *testing.T) {
		s := validStudent()
		s.Email = "invalid-email"
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["email"] != ErrInvalidEmail.Error() {
			t.Errorf("expected %q, got %q", ErrInvalidEmail.Error(), errs["email"])
		}
	})

	t.Run("invalid department returns invalid department error", func(t *testing.T) {
		s := validStudent()
		s.Department = "INVALID"
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["department"] != ErrInvalidDepartment.Error() {
			t.Errorf("expected %q, got %q", ErrInvalidDepartment.Error(), errs["department"])
		}
	})

	t.Run("missing department returns required error", func(t *testing.T) {
		s := validStudent()
		s.Department = ""
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["department"] != ErrDepartmentRequired.Error() {
			t.Errorf("expected %q, got %q", ErrDepartmentRequired.Error(), errs["department"])
		}
	})

	t.Run("invalid semester returns semester error", func(t *testing.T) {
		s := validStudent()
		s.Semester = 9
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["semester"] != ErrSemesterInvalid.Error() {
			t.Errorf("expected %q, got %q", ErrSemesterInvalid.Error(), errs["semester"])
		}
	})

	t.Run("invalid age returns age error", func(t *testing.T) {
		s := validStudent()
		s.Age = 17
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["age"] != ErrAgeInvalid.Error() {
			t.Errorf("expected %q, got %q", ErrAgeInvalid.Error(), errs["age"])
		}
	})
}

func TestValidateAdmin(t *testing.T) {
	validAdmin := func() Admin {
		return Admin{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: "Password123!",
		}
	}

	t.Run("valid admin passes", func(t *testing.T) {
		a := validAdmin()
		if err := Validate.Struct(a); err != nil {
			t.Fatalf("expected valid admin to pass, got error: %v", err)
		}
	})

	t.Run("password too short returns too short error", func(t *testing.T) {
		a := validAdmin()
		a.Password = "Pass1!"
		err := Validate.Struct(a)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["password"] != ErrPasswordTooShort.Error() {
			t.Errorf("expected %q, got %q", ErrPasswordTooShort.Error(), errs["password"])
		}
	})

	t.Run("password too long returns max error", func(t *testing.T) {
		a := validAdmin()
		a.Password = "Aa1!" + strings.Repeat("a", 65)
		err := Validate.Struct(a)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["password"] != ErrPasswordTooLong.Error() {
			t.Errorf("expected %q, got %q", ErrPasswordTooLong.Error(), errs["password"])
		}
	})

	t.Run("password without special char returns invalid password error", func(t *testing.T) {
		a := validAdmin()
		a.Password = "Password123"
		err := Validate.Struct(a)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["password"] != ErrInvalidPassword.Error() {
			t.Errorf("expected %q, got %q", ErrInvalidPassword.Error(), errs["password"])
		}
	})

	t.Run("password with whitespace returns invalid password error", func(t *testing.T) {
		a := validAdmin()
		a.Password = "Pass word123!"
		err := Validate.Struct(a)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["password"] != ErrInvalidPassword.Error() {
			t.Errorf("expected %q, got %q", ErrInvalidPassword.Error(), errs["password"])
		}
	})
}

func TestValidateLoginRequest(t *testing.T) {
	t.Run("valid login request passes", func(t *testing.T) {
		l := LoginRequest{
			Email:    "admin@example.com",
			Password: "Password123!",
		}
		if err := Validate.Struct(l); err != nil {
			t.Fatalf("expected valid login request to pass, got error: %v", err)
		}
	})

	t.Run("missing password returns required error", func(t *testing.T) {
		l := LoginRequest{
			Email:    "admin@example.com",
			Password: "",
		}
		err := Validate.Struct(l)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errs := ValidationErrors(err)
		if errs["password"] != ErrPasswordRequired.Error() {
			t.Errorf("expected %q, got %q", ErrPasswordRequired.Error(), errs["password"])
		}
	})
}

func TestValidationErrors_NonValidatorError(t *testing.T) {
	err := errors.New("some random error")
	errs := ValidationErrors(err)
	if errs["error"] != "some random error" {
		t.Errorf("expected 'some random error', got %q", errs["error"])
	}
}

func TestIgnoredJSONField(t *testing.T) {
	type TestStruct struct {
		IgnoredField string `json:"-" validate:"required"`
	}

	s := TestStruct{}
	err := Validate.Struct(s)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	errs := ValidationErrors(err)
	expectedErrStr := "IgnoredField failed validation on 'required'"
	if errs["IgnoredField"] != expectedErrStr {
		t.Errorf("expected %q, got %q", expectedErrStr, errs["IgnoredField"])
	}
}

func TestDefaultFieldError(t *testing.T) {
	type TestStruct struct {
		UnknownField string `json:"unknown_field" validate:"required"`
	}

	s := TestStruct{}
	err := Validate.Struct(s)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	errs := ValidationErrors(err)
	expectedErrStr := "unknown_field failed validation on 'required'"
	if errs["unknown_field"] != expectedErrStr {
		t.Errorf("expected %q, got %q", expectedErrStr, errs["unknown_field"])
	}
}
