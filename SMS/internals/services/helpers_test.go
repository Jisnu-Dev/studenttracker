package services_test

import (
	"reflect"
	"testing"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
)

func ptr[T any](v T) *T                                      { return &v }
func ptrDept(v models.DepartmentType) *models.DepartmentType { return &v }

func TestBuildPatchQuery(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		patch         models.PatchStudent
		expectedQuery string
		expectedArgs  []interface{}
		expectError   bool
	}{
		{
			name: "single field - name only",
			id:   1,
			patch: models.PatchStudent{
				Name: ptr("test name"),
			},
			expectedQuery: `UPDATE student SET name = $1, "updatedAtUTC" = NOW() WHERE id = $2`,
			expectedArgs:  []interface{}{"test name", int64(1)},
		},
		{
			name: "two fields - name and email",
			id:   5,
			patch: models.PatchStudent{
				Name:  ptr("test name"),
				Email: ptr("test@example.com"),
			},
			expectedQuery: `UPDATE student SET name = $1, email = $2, "updatedAtUTC" = NOW() WHERE id = $3`,
			expectedArgs:  []interface{}{"test name", "test@example.com", int64(5)},
		},
		{
			name: "all fields provided",
			id:   10,
			patch: models.PatchStudent{
				Name:       ptr("test name"),
				Email:      ptr("test@example.com"),
				Department: ptrDept(models.CSE),
				Semester:   ptr(4),
				Age:        ptr(21),
			},
			expectedQuery: `UPDATE student SET name = $1, email = $2, department = $3, semester = $4, age = $5, "updatedAtUTC" = NOW() WHERE id = $6`,
			expectedArgs:  []interface{}{"test name", "test@example.com", models.DepartmentType("CSE"), 4, 21, int64(10)},
		},
		{
			name: "single field - semester only",
			id:   3,
			patch: models.PatchStudent{
				Semester: ptr(6),
			},
			expectedQuery: `UPDATE student SET semester = $1, "updatedAtUTC" = NOW() WHERE id = $2`,
			expectedArgs:  []interface{}{6, int64(3)},
		},
		{
			name: "single field - email only",
			id:   8,
			patch: models.PatchStudent{
				Email: ptr("test@example.com"),
			},
			expectedQuery: `UPDATE student SET email = $1, "updatedAtUTC" = NOW() WHERE id = $2`,
			expectedArgs:  []interface{}{"test@example.com", int64(8)},
		},
		{
			name: "single field - department only",
			id:   12,
			patch: models.PatchStudent{
				Department: ptrDept(models.ECE),
			},
			expectedQuery: `UPDATE student SET department = $1, "updatedAtUTC" = NOW() WHERE id = $2`,
			expectedArgs:  []interface{}{models.DepartmentType("ECE"), int64(12)},
		},
		{
			name: "single field - age only",
			id:   15,
			patch: models.PatchStudent{
				Age: ptr(25),
			},
			expectedQuery: `UPDATE student SET age = $1, "updatedAtUTC" = NOW() WHERE id = $2`,
			expectedArgs:  []interface{}{25, int64(15)},
		},
		{
			name: "two fields - name and department (email skipped)",
			id:   2,
			patch: models.PatchStudent{
				Name:       ptr("Alice"),
				Department: ptrDept(models.IT),
			},
			expectedQuery: `UPDATE student SET name = $1, department = $2, "updatedAtUTC" = NOW() WHERE id = $3`,
			expectedArgs:  []interface{}{"Alice", models.DepartmentType("IT"), int64(2)},
		},
		{
			name: "two fields - name and age (email, department, semester skipped)",
			id:   6,
			patch: models.PatchStudent{
				Name: ptr("Bob"),
				Age:  ptr(30),
			},
			expectedQuery: `UPDATE student SET name = $1, age = $2, "updatedAtUTC" = NOW() WHERE id = $3`,
			expectedArgs:  []interface{}{"Bob", 30, int64(6)},
		},
		{
			name: "two fields - email and semester (name, department skipped)",
			id:   9,
			patch: models.PatchStudent{
				Email:    ptr("carol@example.com"),
				Semester: ptr(3),
			},
			expectedQuery: `UPDATE student SET email = $1, semester = $2, "updatedAtUTC" = NOW() WHERE id = $3`,
			expectedArgs:  []interface{}{"carol@example.com", 3, int64(9)},
		},
		{
			name: "three fields - name, department, and age",
			id:   20,
			patch: models.PatchStudent{
				Name:       ptr("Dana"),
				Department: ptrDept(models.MECH),
				Age:        ptr(22),
			},
			expectedQuery: `UPDATE student SET name = $1, department = $2, age = $3, "updatedAtUTC" = NOW() WHERE id = $4`,
			expectedArgs:  []interface{}{"Dana", models.DepartmentType("MECH"), 22, int64(20)},
		},
		{
			name: "two fields - department and age (semester skipped)",
			id:   11,
			patch: models.PatchStudent{
				Department: ptrDept(models.EEE),
				Age:        ptr(19),
			},
			expectedQuery: `UPDATE student SET department = $1, age = $2, "updatedAtUTC" = NOW() WHERE id = $3`,
			expectedArgs:  []interface{}{models.DepartmentType("EEE"), 19, int64(11)},
		},
		{
			name: "zero-value pointer is still treated as provided - age = 0",
			id:   4,
			patch: models.PatchStudent{
				Age: ptr(0),
			},
			expectedQuery: `UPDATE student SET age = $1, "updatedAtUTC" = NOW() WHERE id = $2`,
			expectedArgs:  []interface{}{0, int64(4)},
		},
		{
			name: "zero-value pointer is still treated as provided - name = empty string",
			id:   7,
			patch: models.PatchStudent{
				Name: ptr(""),
			},
			expectedQuery: `UPDATE student SET name = $1, "updatedAtUTC" = NOW() WHERE id = $2`,
			expectedArgs:  []interface{}{"", int64(7)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, args, err := services.BuildPatchQuery(tt.id, tt.patch)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if query != tt.expectedQuery {
				t.Errorf("expected query:\n  %q\ngot:\n  %q", tt.expectedQuery, query)
			}

			if !reflect.DeepEqual(args, tt.expectedArgs) {
				t.Errorf("expected args %v, got %v", tt.expectedArgs, args)
			}
		})
	}
}
