package hash

import (
	"context"
	"github.com/nqxcode/auth_microservice/pkg/hashing"
)

func (s *service) Check(ctx context.Context, password, hash string) bool {
	return hashing.CheckPasswordHashWithSalt(password, s.salt, hash)
}
