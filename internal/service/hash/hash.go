package hash

import (
	"context"
	"fmt"
	"github.com/nqxcode/auth_microservice/pkg/hashing"
)

func (s *service) Hash(_ context.Context, password string) (string, error) {
	password, err := hashing.HashPasswordWithSalt(password, s.salt)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	return password, nil
}
