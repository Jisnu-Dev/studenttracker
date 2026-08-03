package validation

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	vutils "github.com/Jisnu-Dev/studenttracker/internals/validation/utils"
)

var validDepartments = map[models.DepartmentType]bool{
	models.CSE:   true,
	models.IT:    true,
	models.ECE:   true,
	models.EEE:   true,
	models.MECH:  true,
	models.CIVIL: true,
	models.AIDS:  true,
	models.AIML:  true,
}

// Field Validations
func ValidateName(name string) error {
	if err := vutils.ValidateCleanWhitespace("name", name); err != nil {
		return err
	}
	length := utf8.RuneCountInString(name)
	if length < 2 || length > 100 {
		return errors.New("name must be between 2 and 100 characters")
	}
	if !vutils.NameRegex.MatchString(name) {
		return errors.New("name may only contain letters separated by single spaces")
	}
	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if vutils.IsBlank(email) {
		return errors.New("email cannot be empty or only whitespace")
	}
	if strings.ContainsAny(email, " \t") {
		return errors.New("email cannot contain spaces or tabs")
	}
	if !strings.Contains(email, "@") {
		return errors.New("email must contain '@'")
	}
	if len(email) > vutils.MaxEmailLength {
		return fmt.Errorf("email cannot exceed %d characters", vutils.MaxEmailLength)
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("email is invalid")
	}
	if addr.Address != email {
		return errors.New("email must be a plain address, not a display name")
	}

	return nil
}

func ValidateDepartment(dept models.DepartmentType) error {
	if !validDepartments[dept] {
		return fmt.Errorf("invalid department '%s'; allowed values: CSE, IT, ECE, EEE, MECH, CIVIL, AIDS, AIML", dept)
	}
	return nil
}

func ValidateSemester(semester int) error {
	if semester < 1 || semester > 8 {
		return errors.New("semester must be between 1 and 8")
	}
	return nil
}

func ValidateAge(age int) error {
	if age < 18 || age > 60 {
		return errors.New("age must be between 18 and 60")
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	length := utf8.RuneCountInString(password)
	if length < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if length > vutils.MaxPasswordLength {
		return fmt.Errorf("password cannot exceed %d characters", vutils.MaxPasswordLength)
	}

	var hasLetter, hasDigit bool
	for _, ch := range password {
		switch {
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r':
			return errors.New("password cannot contain whitespace characters")
		case ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z':
			hasLetter = true
		case ch >= '0' && ch <= '9':
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return errors.New("password must contain at least one letter and one number")
	}
	return nil
}

func ValidateAlphaNumeric(fieldName, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	if !vutils.AlphaNumericRegex.MatchString(value) {
		return fmt.Errorf("%s may only contain letters and numbers, no spaces or special characters", fieldName)
	}
	return nil
}

func ValidateStudent(s *models.Student) error {
	errs := vutils.ValidationErrors{}
	errs.Add("name", ValidateName(s.Name))
	errs.Add("email", ValidateEmail(s.Email))
	errs.Add("department", ValidateDepartment(s.Department))
	errs.Add("semester", ValidateSemester(s.Semester))
	errs.Add("age", ValidateAge(s.Age))
	return errs.AsError()
}

func ValidatePatchStudent(p *models.PatchStudent) error {
	errs := vutils.ValidationErrors{}

	if p.Name == nil && p.Email == nil && p.Department == nil && p.Semester == nil && p.Age == nil {
		errs.Add("request", errors.New("at least one field must be provided for update"))
		return errs.AsError()
	}

	if p.Name != nil {
		errs.Add("name", ValidateName(*p.Name))
	}
	if p.Email != nil {
		errs.Add("email", ValidateEmail(*p.Email))
	}
	if p.Department != nil {
		errs.Add("department", ValidateDepartment(*p.Department))
	}
	if p.Semester != nil {
		errs.Add("semester", ValidateSemester(*p.Semester))
	}
	if p.Age != nil {
		errs.Add("age", ValidateAge(*p.Age))
	}

	return errs.AsError()
}

func ValidateAdminRegister(a *models.Admin) error {
	errs := vutils.ValidationErrors{}
	errs.Add("name", ValidateName(a.Name))
	errs.Add("email", ValidateEmail(a.Email))
	errs.Add("password", ValidatePassword(a.Password))
	return errs.AsError()
}

func ValidateAdminLogin(a *models.LoginRequest) error {
	errs := vutils.ValidationErrors{}
	errs.Add("email", ValidateEmail(a.Email))
	errs.Add("password", ValidatePassword(a.Password))
	return errs.AsError()
}
