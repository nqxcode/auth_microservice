package validator

import "github.com/nqxcode/auth_microservice/internal/repository"

type validator struct {
	userRepository repository.UserRepository
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
func NewValidator(userRepository repository.UserRepository) *validator {
	return &validator{
		userRepository: userRepository,
	}
}
