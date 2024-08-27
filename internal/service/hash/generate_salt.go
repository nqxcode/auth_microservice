package hash

import (
	"context"
	"fmt"

	"github.com/nqxcode/auth_microservice/pkg/hashing"
)

// GenerateSalt generate salt
func (s *service) GenerateSalt(_ context.Context) (string, error) {
	salt, err := hashing.GenerateSalt()
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %v", err)
	}

	return salt, nil
}
