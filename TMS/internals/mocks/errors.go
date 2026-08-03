package mocks

type MockOpError int

const (
	OpNone MockOpError = iota
	OpInvalidToken
	OpSigningMethodMismatch
	OpInternalError
)
