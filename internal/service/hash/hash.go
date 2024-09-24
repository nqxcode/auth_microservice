package hash

import (
	"context"
	"fmt"
	"github.com/nqxcode/auth_microservice/internal/utils"
)

// Hash get hash for password
func (s *service) Hash(ctx context.Context, password string) (string, error) {
	hash, err := utils.HashPasswordWithSalt(ctx, password, s.salt)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	return hash, nil
}
