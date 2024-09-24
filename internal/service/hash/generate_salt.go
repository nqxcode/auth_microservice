package hash

import (
	"context"
	"fmt"
	"github.com/nqxcode/auth_microservice/internal/utils"
)

// GenerateSalt generate salt
func (s *service) GenerateSalt(ctx context.Context) (string, error) {
	salt, err := utils.GenerateSalt(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %v", err)
	}

	return salt, nil
}
