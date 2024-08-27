package validator

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/model"
)

// ValidateUser validates the user info
func (v *validator) ValidateUser(ctx context.Context, userInfo model.UserInfo, password, passwordConfirm string) error {
	if userInfo.Name == "" {
		return NewValidationError("name is required")
	}

	if userInfo.Email == "" {
		return NewValidationError("email is required")
	}

	if !ValidateEmail(userInfo.Email) {
		return NewValidationError("invalid email format")
	}

	if password != passwordConfirm {
		return NewValidationError("passwords do not match")
	}

	return nil
}
