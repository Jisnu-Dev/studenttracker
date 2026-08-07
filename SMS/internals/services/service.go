package services

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
)

type ServiceInterface interface {
	CreateStudent(ctx context.Context, student models.Student) (int, error)
	GetAllStudents(ctx context.Context) ([]models.Student, error)
	GetStudentByID(ctx context.Context, id int64) (models.Student, error)
	DeleteStudent(ctx context.Context, id int64) error
	UpdateStudent(ctx context.Context, id int64, student models.Student) error
	PatchStudent(ctx context.Context, id int64, student models.PatchStudent) error
	RegisterAdmin(ctx context.Context, admin models.Admin) (int64, error)
	GetAdminByEmail(ctx context.Context, email string) (models.Admin, error)
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateStudent(ctx context.Context, student models.Student) (int, error) {
	query := `INSERT INTO student (name, email, department, semester, age) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int

	if err := s.db.QueryRowContext(ctx, query, student.Name, student.Email, student.Department, student.Semester, student.Age).Scan(&id); err != nil {
		if IsUniqueViolation(err) {
			return 0, ErrStudentEmailExists
		}
		slog.Error("database query failed",
			slog.String("function", "CreateStudent"),
			slog.String("email", student.Email),
			slog.Any("error", err),
		)
		return 0, ErrCreateStudentFailed
	}

	return id, nil
}

func (s *Service) GetAllStudents(ctx context.Context) ([]models.Student, error) {
	query := `SELECT id, name, email, department, semester, age, "createdAtUTC", "updatedAtUTC" FROM student`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		slog.Error("database query failed",
			slog.String("function", "GetAllStudents"),
			slog.Any("error", err),
		)
		return nil, ErrGetAllStudentsFailed
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
			return nil, ErrGetAllStudentsFailed
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
		return nil, ErrGetAllStudentsFailed
	}

	return students, nil
}

func (s *Service) GetStudentByID(ctx context.Context, id int64) (models.Student, error) {
	query := `SELECT id, name, email, department, semester, age, "createdAtUTC", "updatedAtUTC" FROM student WHERE id = $1`

	var student models.Student
	if err := s.db.QueryRowContext(ctx, query, id).Scan(&student.ID, &student.Name, &student.Email, &student.Department, &student.Semester, &student.Age, &student.CreatedAtUTC, &student.UpdatedAtUTC); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Student{}, nil
		}
		slog.Error("database query failed",
			slog.String("function", "GetStudentByID"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return models.Student{}, ErrGetStudentByIDFailed
	}
	return student, nil
}

func (s *Service) DeleteStudent(ctx context.Context, id int64) error {
	query := `DELETE FROM student WHERE id = $1`

	results, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("database exec failed",
			slog.String("function", "DeleteStudent"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return ErrDeleteStudentFailed
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected",
			slog.String("function", "DeleteStudent"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return ErrDeleteStudentFailed
	}
	if rowsAffected > 1 {
		slog.Error("multiple rows affected",
			slog.String("function", "DeleteStudent"),
			slog.Int64("studentID", id),
			slog.Int64("rowsAffected", rowsAffected),
		)
		return ErrDeleteStudentFailed
	}

	return nil
}

func (s *Service) UpdateStudent(ctx context.Context, id int64, student models.Student) error {
	query := `UPDATE student SET name = $1, email = $2, department = $3, semester = $4, age = $5, "updatedAtUTC" = NOW() WHERE id = $6`

	result, err := s.db.ExecContext(ctx, query, student.Name, student.Email, student.Department, student.Semester, student.Age, id)
	if err != nil {
		if IsUniqueViolation(err) {
			return ErrStudentEmailExists
		}
		slog.Error("database exec failed",
			slog.String("function", "UpdateStudent"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return ErrUpdateStudentFailed
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected",
			slog.String("function", "UpdateStudent"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return ErrUpdateStudentFailed
	}
	if rowsAffected > 1 {
		slog.Error("multiple rows affected",
			slog.String("function", "UpdateStudent"),
			slog.Int64("studentID", id),
			slog.Int64("rowsAffected", rowsAffected),
		)
		return ErrUpdateStudentFailed
	}
	return nil
}

func (s *Service) PatchStudent(ctx context.Context, id int64, student models.PatchStudent) error {
	query, args, err := BuildPatchQuery(id, student)
	if err != nil {
		return ErrNoFieldsToUpdate
	}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		if IsUniqueViolation(err) {
			return ErrStudentEmailExists
		}
		slog.Error("database exec failed",
			slog.String("function", "PatchStudent"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return ErrPatchStudentFailed
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected",
			slog.String("function", "PatchStudent"),
			slog.Int64("studentID", id),
			slog.Any("error", err),
		)
		return ErrPatchStudentFailed
	}
	if rowsAffected > 1 {
		slog.Error("multiple rows affected",
			slog.String("function", "PatchStudent"),
			slog.Int64("studentID", id),
			slog.Int64("rowsAffected", rowsAffected),
		)
		return ErrPatchStudentFailed
	}
	return nil
}

// Admin services
func (s *Service) RegisterAdmin(ctx context.Context, admin models.Admin) (int64, error) {
	query := `INSERT INTO admin (name, email, "passwordHash") VALUES ($1, $2, $3) RETURNING "ID"`
	var id int64

	if err := s.db.QueryRowContext(ctx, query, admin.Name, admin.Email, admin.Password).Scan(&id); err != nil {
		if IsUniqueViolation(err) {
			return 0, ErrAdminEmailExists
		}
		slog.Error("database query failed",
			slog.String("function", "RegisterAdmin"),
			slog.String("email", admin.Email),
			slog.Any("error", err),
		)
		return 0, ErrRegisterAdminFailed
	}

	return id, nil
}

func (s *Service) GetAdminByEmail(ctx context.Context, email string) (models.Admin, error) {
	query := `SELECT "ID", name, email, "passwordHash" FROM admin WHERE email = $1`

	var admin models.Admin
	if err := s.db.QueryRowContext(ctx, query, email).Scan(&admin.ID, &admin.Name, &admin.Email, &admin.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Admin{}, ErrAdminNotFound
		}
		slog.Error("database query failed",
			slog.String("function", "GetAdminByEmail"),
			slog.String("email", email),
			slog.Any("error", err),
		)
		return models.Admin{}, ErrGetAdminByEmailFailed
	}

	return admin, nil
}
