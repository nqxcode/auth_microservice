package converter

import "github.com/nqxcode/auth_microservice/internal/model"

// ToRoleString to role string
func ToRoleString(role int32) string {
	switch role {
	case 1:
		return model.AdminRole
	case 2:
		return model.UserRole
	default:
		return model.UnknownRole
	}
}

// ToRole to role
func ToRole(role string) int32 {
	switch role {
	case model.AdminRole:
		return 1
	case model.UserRole:
		return 2
	default:
		return 0
	}
}
