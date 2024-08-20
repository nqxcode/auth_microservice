package hash

import (
	def "github.com/nqxcode/auth_microservice/internal/service"
)

type service struct {
	salt string
}

// NewService new hash service
func NewService(salt string) def.HashService {
	return &service{
		salt: salt,
	}
}
