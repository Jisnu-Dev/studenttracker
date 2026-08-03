package services

import (
	"database/sql"
	"errors"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	serviceErrors "github.com/Jisnu-Dev/studenttracker/internals/services/errors"
	"github.com/Jisnu-Dev/studenttracker/internals/services/utils"
	"github.com/lib/pq"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}

func (s *Service) CreateStudent(student models.Student) (int, error) {
	query := `INSERT INTO student (name, email, department, semester, age) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int

	if err := s.db.QueryRow(query, student.Name, student.Email, student.Department, student.Semester, student.Age).Scan(&id); err != nil {
		if isUniqueViolation(err) {
			return 0, serviceErrors.ErrStudentEmailExists
		}
		return 0, serviceErrors.ErrCreateStudentFailed
	}

	return id, nil
}

func (s *Service) GetAllStudents() ([]models.Student, error) {
	query := `SELECT id, name, email, department, semester, age, "createdAtUTC", "updatedAtUTC" FROM student`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, serviceErrors.ErrGetAllStudentsFailed
	}
	defer rows.Close()

	students := make([]models.Student, 0)
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Department, &student.Semester, &student.Age, &student.CreatedAtUTC, &student.UpdatedAtUTC); err != nil {
			return nil, serviceErrors.ErrGetAllStudentsFailed
		}

		student.CreatedAtUTC = student.CreatedAtUTC.UTC()
		student.UpdatedAtUTC = student.UpdatedAtUTC.UTC()
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, serviceErrors.ErrGetAllStudentsFailed
	}

	return students, nil
}

func (s *Service) GetStudentByID(id int64) (models.Student, error) {
	query := `SELECT id, name, email, department, semester, age, "createdAtUTC", "updatedAtUTC" FROM student WHERE id = $1`

	var student models.Student
	if err := s.db.QueryRow(query, id).Scan(&student.ID, &student.Name, &student.Email, &student.Department, &student.Semester, &student.Age, &student.CreatedAtUTC, &student.UpdatedAtUTC); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Student{}, serviceErrors.ErrStudentNotFound
		}
		return models.Student{}, serviceErrors.ErrGetStudentByIDFailed
	}
	return student, nil
}

func (s *Service) DeleteStudent(id int64) error {
	query := `DELETE FROM student WHERE id = $1`

	results, err := s.db.Exec(query, id)
	if err != nil {
		return serviceErrors.ErrDeleteStudentFailed
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return serviceErrors.ErrDeleteStudentFailed
	}
	if rowsAffected == 0 {
		return serviceErrors.ErrStudentNotFound
	}

	return nil
}

func (s *Service) UpdateStudent(id int64, student models.Student) error {
	query := `UPDATE student SET name = $1, email = $2, department = $3, semester = $4, age = $5, "updatedAtUTC" = NOW() WHERE id = $6`

	result, err := s.db.Exec(query, student.Name, student.Email, student.Department, student.Semester, student.Age, id)
	if err != nil {
		if isUniqueViolation(err) {
			return serviceErrors.ErrStudentEmailExists
		}
		return serviceErrors.ErrUpdateStudentFailed
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return serviceErrors.ErrUpdateStudentFailed
	}
	if rowsAffected == 0 {
		return serviceErrors.ErrStudentNotFound
	}
	return nil
}

func (s *Service) PatchStudent(id int64, student models.PatchStudent) error {
	query, args, err := utils.BuildPatchQuery(id, student)
	if err != nil {
		return serviceErrors.ErrNoFieldsToUpdate
	}

	result, err := s.db.Exec(query, args...)
	if err != nil {
		if isUniqueViolation(err) {
			return serviceErrors.ErrStudentEmailExists
		}
		return serviceErrors.ErrPatchStudentFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return serviceErrors.ErrPatchStudentFailed
	}
	if rowsAffected == 0 {
		return serviceErrors.ErrStudentNotFound
	}
	return nil
}

// Admin services
func (s *Service) RegisterAdmin(admin models.Admin) (int64, error) {
	query := `INSERT INTO admin (name, email, passwordHash) VALUES ($1, $2, $3) RETURNING id`
	var id int64

	if err := s.db.QueryRow(query, admin.Name, admin.Email, admin.Password).Scan(&id); err != nil {
		if isUniqueViolation(err) {
			return 0, serviceErrors.ErrAdminEmailExists
		}
		return 0, serviceErrors.ErrRegisterAdminFailed
	}

	return id, nil
}

func (s *Service) GetAdminByEmail(email string) (models.Admin, error) {
	query := `SELECT id, name, email, passwordHash FROM admin WHERE email = $1`

	var admin models.Admin
	if err := s.db.QueryRow(query, email).Scan(&admin.ID, &admin.Name, &admin.Email, &admin.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Admin{}, serviceErrors.ErrAdminNotFound
		}
		return models.Admin{}, serviceErrors.ErrGetAdminByEmailFailed
	}

	return admin, nil
}
