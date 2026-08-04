package validation_test

import (
	"strings"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/validation"
)

func strPtr(s string) *string                       { return &s }
func intPtr(i int) *int                             { return &i }
func deptPtr(d models.DepartmentType) *models.DepartmentType { return &d }

func checkErr(t *testing.T, name string, err error, wantErr string) {
	t.Helper()
	if wantErr == "" && err != nil {
		t.Errorf("%s: unexpected error: %v", name, err)
	} else if wantErr != "" && (err == nil || !strings.Contains(err.Error(), wantErr)) {
		t.Errorf("%s: expected error containing %q, got %v", name, wantErr, err)
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct{ name, in, wantErr string }{
		{"valid", "Alice Smith", ""},
		{"unicode", "Ñoño García", ""},
		{"empty", "", "name is required"},
		{"blank", "   ", "name cannot be empty or only whitespace"},
		{"leading space", " Alice", "name cannot have leading or trailing spaces"},
		{"tab", "Alice\tSmith", "name cannot contain tab characters"},
		{"multiple spaces", "Alice  Smith", "name cannot contain multiple consecutive spaces"},
		{"too short", "A", "name must be between 2 and 100 characters"},
		{"too long", strings.Repeat("A", 101), "name must be between 2 and 100 characters"},
		{"digits/special", "Alice123", "name may only contain letters separated by single spaces"},
	}
	for _, tt := range tests {
		checkErr(t, tt.name, validation.ValidateName(tt.in), tt.wantErr)
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct{ name, in, wantErr string }{
		{"valid", "user@example.com", ""},
		{"plus tag", "user+tag@example.com", ""},
		{"empty", "", "email is required"},
		{"blank", "   ", "email cannot be empty or only whitespace"},
		{"spaces", "us er@example.com", "email cannot contain spaces or tabs"},
		{"no at", "userexample.com", "email must contain '@'"},
		{"too long", strings.Repeat("a", 250) + "@b.com", "email cannot exceed 254 characters"},
		{"invalid format", "not-an-email@", "email is invalid"},
	}
	for _, tt := range tests {
		checkErr(t, tt.name, validation.ValidateEmail(tt.in), tt.wantErr)
	}
}

func TestValidateDepartment(t *testing.T) {
	for _, d := range []models.DepartmentType{models.CSE, models.IT, models.ECE, models.EEE, models.MECH, models.CIVIL, models.AIDS, models.AIML} {
		checkErr(t, string(d), validation.ValidateDepartment(d), "")
	}
	checkErr(t, "invalid", validation.ValidateDepartment("UNKNOWN"), "invalid department")
	checkErr(t, "empty", validation.ValidateDepartment(""), "invalid department")
}

func TestValidateSemesterAndAge(t *testing.T) {
	semTests := []struct {
		sem     int
		wantErr string
	}{
		{1, ""}, {8, ""}, {0, "between 1 and 8"}, {9, "between 1 and 8"},
	}
	for _, tt := range semTests {
		checkErr(t, "semester", validation.ValidateSemester(tt.sem), tt.wantErr)
	}

	ageTests := []struct {
		age     int
		wantErr string
	}{
		{18, ""}, {60, ""}, {17, "between 18 and 60"}, {61, "between 18 and 60"},
	}
	for _, tt := range ageTests {
		checkErr(t, "age", validation.ValidateAge(tt.age), tt.wantErr)
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct{ name, in, wantErr string }{
		{"valid", "Password1!", ""},
		{"empty", "", "password is required"},
		{"too short", "Pass1", "password must be at least 8 characters long"},
		{"too long", strings.Repeat("A1", 37), "password cannot exceed 72 characters"},
		{"whitespace", "Pass word1", "password cannot contain whitespace characters"},
		{"no digit", "OnlyLetters", "password must contain at least one letter and one number"},
		{"no letter", "12345678", "password must contain at least one letter and one number"},
	}
	for _, tt := range tests {
		checkErr(t, tt.name, validation.ValidatePassword(tt.in), tt.wantErr)
	}
}

func TestValidateAlphaNumeric(t *testing.T) {
	tests := []struct{ name, in, wantErr string }{
		{"valid", "abc123XYZ", ""},
		{"empty", "", "token is required"},
		{"space/special", "abc 123", "token may only contain letters and numbers"},
	}
	for _, tt := range tests {
		checkErr(t, tt.name, validation.ValidateAlphaNumeric("token", tt.in), tt.wantErr)
	}
}

func TestValidateStudent(t *testing.T) {
	valid := models.Student{Name: "Alice Smith", Email: "alice@example.com", Department: models.CSE, Semester: 3, Age: 21}
	checkErr(t, "valid", validation.ValidateStudent(&valid), "")

	invalid := valid
	invalid.Name, invalid.Semester = "", 10
	checkErr(t, "invalid", validation.ValidateStudent(&invalid), "name: name is required")
}

func TestValidatePatchStudent(t *testing.T) {
	tests := []struct {
		name    string
		p       models.PatchStudent
		wantErr string
	}{
		{"valid name", models.PatchStudent{Name: strPtr("Bob")}, ""},
		{"valid multi", models.PatchStudent{Department: deptPtr(models.IT), Age: intPtr(22)}, ""},
		{"empty", models.PatchStudent{}, "at least one field must be provided"},
		{"invalid field", models.PatchStudent{Age: intPtr(10)}, "age: age must be between 18 and 60"},
	}
	for _, tt := range tests {
		checkErr(t, tt.name, validation.ValidatePatchStudent(&tt.p), tt.wantErr)
	}
}

func TestValidateAdmin(t *testing.T) {
	reg := models.Admin{Name: "Admin", Email: "admin@example.com", Password: "Password1"}
	checkErr(t, "valid reg", validation.ValidateAdminRegister(&reg), "")

	badReg := models.Admin{Name: "", Email: "admin@example.com", Password: "123"}
	checkErr(t, "invalid reg", validation.ValidateAdminRegister(&badReg), "name: name is required")

	login := models.LoginRequest{Email: "admin@example.com", Password: "Password1"}
	checkErr(t, "valid login", validation.ValidateAdminLogin(&login), "")

	badLogin := models.LoginRequest{Email: "bad", Password: ""}
	checkErr(t, "invalid login", validation.ValidateAdminLogin(&badLogin), "password: password is required")
}
