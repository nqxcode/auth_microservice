package converter

import (
	"github.com/nqxcode/auth_microservice/internal/model"
	modelRepo "github.com/nqxcode/auth_microservice/internal/repository/user/pg/model"
	"github.com/samber/lo"
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

// ToManyUserFromRepo convert to many user models
func ToManyUserFromRepo(users []modelRepo.User) []model.User {
	return lo.Map(users, func(user modelRepo.User, _ int) model.User {
		return *ToUserFromRepo(&user)
	})
}
