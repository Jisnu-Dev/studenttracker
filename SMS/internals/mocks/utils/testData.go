package utils

import "github.com/Jisnu-Dev/studenttracker/internals/models"

// ValidStudent returns a models.Student struct that satisfies all SMS validation rules.
// Useful for captured argument assertions in handler tests.
func ValidStudent() models.Student {
	return models.Student{
		Name:       "test name",
		Email:      "test@example.com",
		Department: models.CSE,
		Semester:   3,
		Age:        20,
	}
}

// ValidStudentJSON is the JSON body string representation of ValidStudent.
// Use this as the raw body string in handler test tables.
const ValidStudentJSON = `{"name":"test name","email":"test@example.com","department":"CSE","semester":3,"age":20}`
