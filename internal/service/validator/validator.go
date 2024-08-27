package validator

type validator struct {
}

// ValidationError validation error
type ValidationError struct {
	message string
}

// NewValidationError new validation error
func NewValidationError(message string) *ValidationError {
	return &ValidationError{message: message}
}

// Error get error message
func (e *ValidationError) Error() string {
	return e.message
}

// NewValidator new validator
func NewValidator() *validator {
	return &validator{}
}
