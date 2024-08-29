package converter

import (
	"github.com/nqxcode/platform_common/pointer"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nqxcode/auth_microservice/internal/model"
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
	return lo.Map(users, func(user model.User, _ int) *desc.User {
		return ToUserFromService(&user)
	})
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
		name = pointer.ToPtr(info.GetName().GetValue())
	}

	var role *int32
	if info.GetRole() != 0 {
		role = lo.ToPtr(int32(info.GetRole()))
	}

	return &model.UpdateUserInfo{
		Name: name,
		Role: role,
	}
}

// ToUserInfoFromDesc to user model
func ToUserInfoFromDesc(info *desc.UserInfo) *model.UserInfo {
	var (
		name, email string
		role        int32
	)

	if info != nil {
		name = info.GetName()
		email = info.GetEmail()
		role = int32(info.GetRole())
	}

	return &model.UserInfo{
		Name:  name,
		Email: email,
		Role:  role,
	}
}

func ToUserInfoFromMessage(info *model.UserInfoInMessage) *model.UserInfo {
	return &model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  info.Role,
	}
}
