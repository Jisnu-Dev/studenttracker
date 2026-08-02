package services

import (
	"database/sql"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services/utils"
)

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
		return 0, err
	}

	return id, nil
}

func (s *Service) GetAllStudents() ([]models.Student, error) {
	query := `SELECT id, name, email, department, semester, age, "createdAtUTC", "updatedAtUTC" FROM student`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := make([]models.Student, 0)
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Department, &student.Semester, &student.Age, &student.CreatedAtUTC, &student.UpdatedAtUTC); err != nil {
			return nil, err
		}

		student.CreatedAtUTC = student.CreatedAtUTC.UTC()
		student.UpdatedAtUTC = student.UpdatedAtUTC.UTC()
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (s *Service) GetStudentByID(id int64) (models.Student, error) {
	query := `SELECT id, name, email, department, semester, age, "createdAtUTC", "updatedAtUTC" FROM student WHERE id = $1`

	var student models.Student
	if err := s.db.QueryRow(query, id).Scan(&student.ID, &student.Name, &student.Email, &student.Department, &student.Semester, &student.Age, &student.CreatedAtUTC, &student.UpdatedAtUTC); err != nil {
		return models.Student{}, err
	}
	return student, nil
}

func (s *Service) DeleteStudent(id int64) error {
	query := `DELETE FROM student WHERE id = $1`

	results, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Service) UpdateStudent(id int64, student models.Student) error {
	query := `UPDATE student SET name = $1, email = $2, department = $3, semester = $4, age = $5, "updatedAtUTC" = NOW() WHERE id = $6`

	result, err := s.db.Exec(query, student.Name, student.Email, student.Department, student.Semester, student.Age, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Service) PatchStudent(id int64, student models.PatchStudent) error {
	query, args, err := utils.BuildPatchQuery(id, student)
	if err != nil {
		return err
	}

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Service) RegisterAdmin(admin models.Admin) (int64, error) {
	query := `INSERT INTO admin (name, email, passwordHash) VALUES ($1, $2, $3) RETURNING id`
	var id int64

	if err := s.db.QueryRow(query, admin.Name, admin.Email, admin.Password).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
