package mocks

import (
	"github.com/Jisnu-Dev/studenttracker/internals/models"
	serviceErrors "github.com/Jisnu-Dev/studenttracker/internals/services/errors"
)

type MockService struct {
	// 1. Error simulation triggers
	CreateStudentError   MockOpError
	GetAllStudentsError  MockOpError
	GetStudentByIDError  MockOpError
	DeleteStudentError   MockOpError
	UpdateStudentError   MockOpError
	PatchStudentError    MockOpError
	RegisterAdminError   MockOpError
	GetAdminByEmailError MockOpError

	// 2. Mock return data configuration
	ReturnEmptyStudents bool
	Students            []models.Student
	Student             models.Student
	Admin               models.Admin
	CreatedStudentID    int
	RegisteredAdminID   int64

	// 3. Captured arguments (for verification in tests)
	CapturedID           int64
	CapturedEmail        string
	CapturedStudent      models.Student
	CapturedPatchStudent models.PatchStudent
	CapturedAdmin        models.Admin
}

func (m *MockService) CreateStudent(student models.Student) (int, error) {
	m.CapturedStudent = student

	switch m.CreateStudentError {
	case OpEmailExists:
		return 0, serviceErrors.ErrStudentEmailExists
	case OpInternalError:
		return 0, serviceErrors.ErrCreateStudentFailed
	}

	if m.CreatedStudentID != 0 {
		return m.CreatedStudentID, nil
	}
	return 1, nil
}

func (m *MockService) GetAllStudents() ([]models.Student, error) {
	switch m.GetAllStudentsError {
	case OpInternalError:
		return nil, serviceErrors.ErrGetAllStudentsFailed
	}

	if m.ReturnEmptyStudents {
		return []models.Student{}, nil
	}

	if len(m.Students) > 0 {
		return m.Students, nil
	}

	return []models.Student{
		{ID: 1, Name: "test", Email: "test@example.com", Department: "CS", Semester: 4, Age: 21},
	}, nil
}

func (m *MockService) GetStudentByID(id int64) (models.Student, error) {
	m.CapturedID = id

	switch m.GetStudentByIDError {
	case OpNotFound:
		return models.Student{}, serviceErrors.ErrStudentNotFound
	case OpInternalError:
		return models.Student{}, serviceErrors.ErrGetStudentByIDFailed
	}

	if m.Student.ID != 0 {
		return m.Student, nil
	}
	return models.Student{ID: id, Name: "test", Email: "test@example.com"}, nil
}

func (m *MockService) DeleteStudent(id int64) error {
	m.CapturedID = id

	switch m.DeleteStudentError {
	case OpNotFound:
		return serviceErrors.ErrStudentNotFound
	case OpInternalError:
		return serviceErrors.ErrDeleteStudentFailed
	}

	return nil
}

func (m *MockService) UpdateStudent(id int64, student models.Student) error {
	m.CapturedID = id
	m.CapturedStudent = student

	switch m.UpdateStudentError {
	case OpNotFound:
		return serviceErrors.ErrStudentNotFound
	case OpEmailExists:
		return serviceErrors.ErrStudentEmailExists
	case OpInternalError:
		return serviceErrors.ErrUpdateStudentFailed
	}

	return nil
}

func (m *MockService) PatchStudent(id int64, student models.PatchStudent) error {
	m.CapturedID = id
	m.CapturedPatchStudent = student

	switch m.PatchStudentError {
	case OpNotFound:
		return serviceErrors.ErrStudentNotFound
	case OpEmailExists:
		return serviceErrors.ErrStudentEmailExists
	case OpNoFieldsToUpdate:
		return serviceErrors.ErrNoFieldsToUpdate
	case OpInternalError:
		return serviceErrors.ErrPatchStudentFailed
	}

	return nil
}

func (m *MockService) RegisterAdmin(admin models.Admin) (int64, error) {
	m.CapturedAdmin = admin

	switch m.RegisterAdminError {
	case OpEmailExists:
		return 0, serviceErrors.ErrAdminEmailExists
	case OpInternalError:
		return 0, serviceErrors.ErrRegisterAdminFailed
	}

	if m.RegisteredAdminID != 0 {
		return m.RegisteredAdminID, nil
	}
	return 1, nil
}

func (m *MockService) GetAdminByEmail(email string) (models.Admin, error) {
	m.CapturedEmail = email

	switch m.GetAdminByEmailError {
	case OpNotFound:
		return models.Admin{}, serviceErrors.ErrAdminNotFound
	case OpInternalError:
		return models.Admin{}, serviceErrors.ErrGetAdminByEmailFailed
	}

	if m.Admin.Email != "" {
		return m.Admin, nil
	}
	return models.Admin{ID: 1, Name: "Admin", Email: email, Password: "$2a$10$dummyhashedpassword"}, nil
}
