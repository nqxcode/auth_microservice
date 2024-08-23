package converter

import (
	"github.com/nqxcode/auth_microservice/internal/model"
	modelRepo "github.com/nqxcode/auth_microservice/internal/repository/user/pg/model"
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
	result := make([]model.User, 0, len(users))
	for i := range users {
		m := ToUserFromRepo(&users[i])
		if m != nil {
			result = append(result, *m)
		}
	}

	return result
}
