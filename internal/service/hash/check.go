package hash

import (
	"context"

	"github.com/nqxcode/auth_microservice/internal/utils"
)

// Check password
func (s *service) Check(ctx context.Context, password, hash string) bool {
	return utils.VerifyPassword(ctx, password, s.salt, hash)
}
