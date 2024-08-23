package converter

import (
	"database/sql"
	"time"

	"github.com/nqxcode/auth_microservice/internal/model"
	modelRepo "github.com/nqxcode/auth_microservice/internal/repository/user/redis/model"
)

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
