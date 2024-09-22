package hash

import (
	"context"
	"fmt"

	"github.com/nqxcode/auth_microservice/pkg/hashing"
)

// Hash get hash for password
func (s *service) Hash(_ context.Context, password string) (string, error) {
	hash, err := hashing.HashPasswordWithSalt(password, s.salt)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	return hash, nil
}
