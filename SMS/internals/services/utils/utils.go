package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
)

func BuildPatchQuery(id int64, student models.PatchStudent) (string, []any, error) {
	parts := []string{}
	args := []any{}
	paramIndex := 1

	if student.Name != nil {
		parts = append(parts, "name = $"+strconv.Itoa(paramIndex))
		args = append(args, *student.Name)
		paramIndex++
	}

	if student.Email != nil {
		parts = append(parts, "email = $"+strconv.Itoa(paramIndex))
		args = append(args, *student.Email)
		paramIndex++
	}

	if student.Department != nil {
		parts = append(parts, "department = $"+strconv.Itoa(paramIndex))
		args = append(args, *student.Department)
		paramIndex++
	}

	if student.Semester != nil {
		parts = append(parts, "semester = $"+strconv.Itoa(paramIndex))
		args = append(args, *student.Semester)
		paramIndex++
	}

	if student.Age != nil {
		parts = append(parts, "age = $"+strconv.Itoa(paramIndex))
		args = append(args, *student.Age)
		paramIndex++
	}

	if len(parts) == 0 {
		return "", nil, errors.New("no fields to update")
	}

	args = append(args, id)
	query := `UPDATE student SET ` + strings.Join(parts, ", ") + `, "updatedAtUTC" = NOW() WHERE id = $` + strconv.Itoa(paramIndex)

	return query, args, nil
}
