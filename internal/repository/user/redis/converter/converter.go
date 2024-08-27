package converter

import (
	"database/sql"
	"github.com/samber/lo"
	"time"

	"github.com/nqxcode/auth_microservice/internal/model"
	modelRepo "github.com/nqxcode/auth_microservice/internal/repository/user/redis/model"
)

// ToManyUserFromRepo convert to user model
func ToManyUserFromRepo(users []modelRepo.User) []model.User {
	return lo.Map(users, func(user modelRepo.User, _ int) model.User {
		return *ToUserFromRepo(&user)
	})
}

// ToUserFromRepo convert to user model
func ToUserFromRepo(user *modelRepo.User) *model.User {
	var updatedAt sql.NullTime
	if user.UpdatedAtNs != nil {
		updatedAt = sql.NullTime{
			Time:  time.Unix(0, *user.UpdatedAtNs),
			Valid: true,
		}
	}

	return &model.User{
		ID: user.ID,
		Info: model.UserInfo{
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		Password:  user.Password,
		CreatedAt: time.Unix(0, user.CreatedAtNs),
		UpdatedAt: updatedAt,
	}
}
