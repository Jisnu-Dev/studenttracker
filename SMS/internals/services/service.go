package services

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	services "github.com/Jisnu-Dev/studenttracker/internals/services/errors"
	"github.com/Jisnu-Dev/studenttracker/internals/services/utils"
)

type ServiceInterface interface {
	CreateStudent(student models.Student) (int, error)
	GetAllStudents() ([]models.Student, error)
	GetStudentByID(id int64) (models.Student, error)
	DeleteStudent(id int64) error
	UpdateStudent(id int64, student models.Student) error
	PatchStudent(id int64, student models.PatchStudent) error
	RegisterAdmin(admin models.Admin) (int64, error)
	GetAdminByEmail(email string) (models.Admin, error)
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateStudent(student models.Student) (int, error) {
	query := `INSERT INTO student (name, email, department, semester, age) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int

	if err := s.db.QueryRow(query, student.Name, student.Email, student.Department, student.Semester, student.Age).Scan(&id); err != nil {
		if utils.IsUniqueViolation(err) {
			return 0, services.ErrStudentEmailExists
		}
		slog.Error("database query failed",
			slog.String("function", "CreateStudent"),
			slog.String("email", student.Email),
			slog.Any("error", err),
		)
		return 0, services.ErrCreateStudentFailed
	}

	return id, nil
}

func (s *Service) GetAllStudents() ([]models.Student, error) {
	query := `SELECT id, name, email, department, semester, age, "createdAtUTC", "updatedAtUTC" FROM student`

	rows, err := s.db.Query(query)
	if err != nil {
		slog.Error("database query failed",
			slog.String("function", "GetAllStudents"),
			slog.Any("error", err),
		)
		return nil, services.ErrGetAllStudentsFailed
	}
	defer rows.Close()

	students := make([]models.Student, 0)
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Department, &student.Semester, &student.Age, &student.CreatedAtUTC, &student.UpdatedAtUTC); err != nil {
			slog.Error("failed to scan student row",
				slog.String("function", "GetAllStudents"),
				slog.Any("error", err),
			)
			return nil, services.ErrGetAllStudentsFailed
		}

		student.CreatedAtUTC = student.CreatedAtUTC.UTC()
		student.UpdatedAtUTC = student.UpdatedAtUTC.UTC()
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		slog.Error("rows iteration error",
			slog.String("function", "GetAllStudents"),
			slog.Any("error", err),
		)
		return nil, services.ErrGetAllStudentsFailed
	}

	return students, nil
}

func (s *Service) GetStudentByID(id int64) (models.Student, error) {
	query := `SELECT id, name, email, department, semester, age, "createdAtUTC", "updatedAtUTC" FROM student WHERE id = $1`

	var student models.Student
	if err := s.db.QueryRow(query, id).Scan(&student.ID, &student.Name, &student.Email, &student.Department, &student.Semester, &student.Age, &student.CreatedAtUTC, &student.UpdatedAtUTC); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Student{}, services.ErrStudentNotFound
		}
		slog.Error("database query failed",
			slog.String("function", "GetStudentByID"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return models.Student{}, services.ErrGetStudentByIDFailed
	}
	return student, nil
}

func (s *Service) DeleteStudent(id int64) error {
	query := `DELETE FROM student WHERE id = $1`

	results, err := s.db.Exec(query, id)
	if err != nil {
		slog.Error("database exec failed",
			slog.String("function", "DeleteStudent"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return services.ErrDeleteStudentFailed
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected",
			slog.String("function", "DeleteStudent"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return services.ErrDeleteStudentFailed
	}
	if rowsAffected == 0 {
		return services.ErrStudentNotFound
	}

	return nil
}

func (s *Service) UpdateStudent(id int64, student models.Student) error {
	query := `UPDATE student SET name = $1, email = $2, department = $3, semester = $4, age = $5, "updatedAtUTC" = NOW() WHERE id = $6`

	result, err := s.db.Exec(query, student.Name, student.Email, student.Department, student.Semester, student.Age, id)
	if err != nil {
		if utils.IsUniqueViolation(err) {
			return services.ErrStudentEmailExists
		}
		slog.Error("database exec failed",
			slog.String("function", "UpdateStudent"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return services.ErrUpdateStudentFailed
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected",
			slog.String("function", "UpdateStudent"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return services.ErrUpdateStudentFailed
	}
	if rowsAffected == 0 {
		return services.ErrStudentNotFound
	}
	return nil
}

func (s *Service) PatchStudent(id int64, student models.PatchStudent) error {
	query, args, err := utils.BuildPatchQuery(id, student)
	if err != nil {
		return services.ErrNoFieldsToUpdate
	}

	result, err := s.db.Exec(query, args...)
	if err != nil {
		if utils.IsUniqueViolation(err) {
			return services.ErrStudentEmailExists
		}
		slog.Error("database exec failed",
			slog.String("function", "PatchStudent"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return services.ErrPatchStudentFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected",
			slog.String("function", "PatchStudent"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return services.ErrPatchStudentFailed
	}
	if rowsAffected == 0 {
		return services.ErrStudentNotFound
	}
	return nil
}

// Admin services
func (s *Service) RegisterAdmin(admin models.Admin) (int64, error) {
	query := `INSERT INTO admin (name, email, passwordHash) VALUES ($1, $2, $3) RETURNING id`
	var id int64

	if err := s.db.QueryRow(query, admin.Name, admin.Email, admin.Password).Scan(&id); err != nil {
		if utils.IsUniqueViolation(err) {
			return 0, services.ErrAdminEmailExists
		}
		slog.Error("database query failed",
			slog.String("function", "RegisterAdmin"),
			slog.String("email", admin.Email),
			slog.Any("error", err),
		)
		return 0, services.ErrRegisterAdminFailed
	}

	return id, nil
}

func (s *Service) GetAdminByEmail(email string) (models.Admin, error) {
	query := `SELECT id, name, email, passwordHash FROM admin WHERE email = $1`

	var admin models.Admin
	if err := s.db.QueryRow(query, email).Scan(&admin.ID, &admin.Name, &admin.Email, &admin.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Admin{}, services.ErrAdminNotFound
		}
		slog.Error("database query failed",
			slog.String("function", "GetAdminByEmail"),
			slog.String("email", email),
			slog.Any("error", err),
		)
		return models.Admin{}, services.ErrGetAdminByEmailFailed
	}

	return admin, nil
}
