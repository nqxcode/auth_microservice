package auth

import (
	"github.com/nqxcode/auth_microservice/internal/service"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

// Implementation chat api implementation
type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation new auth service implementation
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
