package mocks

// MockOpError represents the simulated error condition to return from a mock operation.
type MockOpError int

const (
	// OpNone indicates no error simulation (mock returns successful default or configured data).
	OpNone MockOpError = iota

	// OpNotFound simulates a resource not found condition (e.g. ErrStudentNotFound, ErrAdminNotFound).
	OpNotFound

	// OpEmailExists simulates a unique constraint violation when an email is already registered.
	OpEmailExists

	// OpInternalError simulates an unexpected database, server, or gRPC communication failure.
	OpInternalError

	// OpNoFieldsToUpdate simulates an update/patch request when no modifiable fields were provided.
	OpNoFieldsToUpdate

	// OpInvalidCredentials simulates invalid authentication credentials during login.
	OpInvalidCredentials

	// OpInvalidToken simulates an invalid, malformed, or expired JWT token.
	OpInvalidToken
)
