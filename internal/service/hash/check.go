package hash

import (
	"context"
	"github.com/nqxcode/auth_microservice/internal/utils"
)

func (s *service) Check(ctx context.Context, password, hash string) bool {
	return utils.CheckPasswordHashWithSalt(ctx, password, s.salt, hash)
}
