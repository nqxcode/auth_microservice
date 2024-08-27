package validator

type validator struct {
}

type ValidationError struct {
	message string
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{message: message}
}

func (e *ValidationError) Error() string {
	return e.message
}

func NewValidator() *validator {
	return &validator{}
}
