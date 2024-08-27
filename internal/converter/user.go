package converter

import (
	"github.com/nqxcode/auth_microservice/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

// ToUserFromService convert to user model
func ToUserFromService(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoFromService(user),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// ToUsersFromService convert to users model
func ToUsersFromService(users []model.User) []*desc.User {
	userList := make([]*desc.User, 0, len(users))
	for i := range users {
		userList = append(userList, ToUserFromService(&users[i]))
	}
	return userList
}

// ToUserInfoFromService convert to user info model
func ToUserInfoFromService(user *model.User) *desc.UserInfo {
	return &desc.UserInfo{
		Name:  user.Info.Name,
		Email: user.Info.Email,
		Role:  desc.Role(user.Info.Role),
	}
}

// ToUpdateUserInfoFromDesc to user info model
func ToUpdateUserInfoFromDesc(info *desc.UpdateUserInfo) *model.UpdateUserInfo {
	var name *string
	if info.GetName() != nil {
		name = toPtr(info.GetName().GetValue())
	}

	var role *int32
	if info.GetRole() != 0 {
		role = toPtr(int32(info.GetRole()))
	}

	return &model.UpdateUserInfo{
		Name: name,
		Role: role,
	}
}

// ToUserFromDesc to user model
func ToUserFromDesc(info *desc.UserInfo, password, passwordConfirm string) *model.User {
	var (
		name, email string
		role        int32
	)

	if info != nil {
		name = info.GetName()
		email = info.GetEmail()
		role = int32(info.GetRole())
	}

	return &model.User{
		Password:        password,
		PasswordConfirm: passwordConfirm,
		Info: model.UserInfo{
			Name:  name,
			Email: email,
			Role:  role,
		},
	}
}

func toPtr[T any](s T) *T {
	return &s
}
