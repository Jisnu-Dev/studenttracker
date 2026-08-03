package mocks

// MockOpError represents the simulated error condition to return from a mock operation.
type MockOpError int

const (
	// OpNone indicates no error simulation (mock returns successful default or configured token data).
	OpNone MockOpError = iota

	// OpInvalidToken simulates an invalid, unparseable, or expired JWT token.
	OpInvalidToken

	// OpSigningMethodMismatch simulates a token signed with an unexpected algorithm (e.g. RSA instead of HMAC).
	OpSigningMethodMismatch

	// OpInternalError simulates an internal token generation or validation failure.
	OpInternalError
)
