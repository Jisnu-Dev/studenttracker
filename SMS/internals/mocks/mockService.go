package mocks

import (
	"context"

	"github.com/Jisnu-Dev/studenttracker/internals/models"
	"github.com/Jisnu-Dev/studenttracker/internals/services"
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

func (m *MockService) CreateStudent(ctx context.Context, student models.Student) (int, error) {
	switch m.CreateStudentError {
	case OpEmailExists:
		return 0, services.ErrStudentEmailExists
	case OpInternalError:
		return 0, services.ErrCreateStudentFailed
	}

	if m.CreatedStudentID != 0 {
		return m.CreatedStudentID, nil
	}
	return 1, nil
}

func (m *MockService) GetAllStudents(ctx context.Context) ([]models.Student, error) {
	switch m.GetAllStudentsError {
	case OpInternalError:
		return nil, services.ErrGetAllStudentsFailed
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

func (m *MockService) GetStudentByID(ctx context.Context, id int64) (models.Student, error) {
	switch m.GetStudentByIDError {
	case OpNotFound:
		return models.Student{}, services.ErrStudentNotFound
	case OpInternalError:
		return models.Student{}, services.ErrGetStudentByIDFailed
	}

	if m.Student.ID != 0 {
		return m.Student, nil
	}
	return models.Student{ID: id, Name: "test", Email: "test@example.com"}, nil
}

func (m *MockService) DeleteStudent(ctx context.Context, id int64) error {
	switch m.DeleteStudentError {
	case OpNotFound:
		return services.ErrStudentNotFound
	case OpInternalError:
		return services.ErrDeleteStudentFailed
	}

	return nil
}

func (m *MockService) UpdateStudent(ctx context.Context, id int64, student models.Student) error {
	switch m.UpdateStudentError {
	case OpNotFound:
		return services.ErrStudentNotFound
	case OpEmailExists:
		return services.ErrStudentEmailExists
	case OpInternalError:
		return services.ErrUpdateStudentFailed
	}

	return nil
}

func (m *MockService) PatchStudent(ctx context.Context, id int64, student models.PatchStudent) error {
	switch m.PatchStudentError {
	case OpNotFound:
		return services.ErrStudentNotFound
	case OpEmailExists:
		return services.ErrStudentEmailExists
	case OpNoFieldsToUpdate:
		return services.ErrNoFieldsToUpdate
	case OpInternalError:
		return services.ErrPatchStudentFailed
	}

	return nil
}

func (m *MockService) RegisterAdmin(ctx context.Context, admin models.Admin) (int64, error) {
	switch m.RegisterAdminError {
	case OpEmailExists:
		return 0, services.ErrAdminEmailExists
	case OpInternalError:
		return 0, services.ErrRegisterAdminFailed
	}

	if m.RegisteredAdminID != 0 {
		return m.RegisteredAdminID, nil
	}
	return 1, nil
}

func (m *MockService) GetAdminByEmail(ctx context.Context, email string) (models.Admin, error) {
	switch m.GetAdminByEmailError {
	case OpNotFound:
		return models.Admin{}, services.ErrAdminNotFound
	case OpInternalError:
		return models.Admin{}, services.ErrGetAdminByEmailFailed
	}

	if m.Admin.Email != "" {
		return m.Admin, nil
	}
	return models.Admin{ID: 1, Name: "Admin", Email: email, Password: "$2a$10$dummyhashedpassword"}, nil
}
