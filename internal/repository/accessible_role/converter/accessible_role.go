package converter

import (
	"github.com/samber/lo"

	"github.com/nqxcode/auth_microservice/internal/model"
	modelRepo "github.com/nqxcode/auth_microservice/internal/repository/accessible_role/model"
)

// ToAccessibleRoleFromRepo model converter
func ToAccessibleRoleFromRepo(user *modelRepo.AccessibleRole) *model.AccessibleRole {
	return &model.AccessibleRole{
		ID:              user.ID,
		Role:            user.Role,
		EndpointAddress: user.EndpointAddress,
		CreatedAt:       user.CreatedAt,
	}
}

// ToManyAccessibleRoleFromRepo convert to many accessible role models
func ToManyAccessibleRoleFromRepo(users []modelRepo.AccessibleRole) []model.AccessibleRole {
	return lo.Map(users, func(user modelRepo.AccessibleRole, _ int) model.AccessibleRole {
		return *ToAccessibleRoleFromRepo(&user)
	})
}
