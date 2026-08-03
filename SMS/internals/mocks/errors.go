package mocks

type MockOpError int

const (
	OpNone MockOpError = iota
	OpNotFound
	OpEmailExists
	OpInternalError
	OpNoFieldsToUpdate
	OpInvalidCredentials
)
