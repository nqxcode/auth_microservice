package converter

import (
	"github.com/nqxcode/auth_microservice/internal/model"
	modelRepo "github.com/nqxcode/auth_microservice/internal/repository/user/model"
)

// ToUserFromRepo model converter
func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID: user.ID,
		Info: model.UserInfo{
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
