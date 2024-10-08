package converter

import (
	"database/sql"
	"time"

	"github.com/nqxcode/platform_common/pointer"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nqxcode/auth_microservice/internal/model"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
)

// ToUserFromService convert to user model
func ToUserFromService(user *model.User) *desc.User {
	if user == nil {
		return nil
	}

	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoFromService(&user.Info),
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
func ToUserInfoFromService(info *model.UserInfo) *desc.UserInfo {
	if info == nil {
		return nil
	}

	return &desc.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  desc.Role(info.Role),
	}
}

// ToUpdateUserInfoFromDesc to user info model
func ToUpdateUserInfoFromDesc(info *desc.UpdateUserInfo) *model.UpdateUserInfo {
	if info == nil {
		return nil
	}

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
	if info == nil {
		return nil
	}

	return &model.UserInfo{
		Name:  info.GetName(),
		Email: info.GetEmail(),
		Role:  int32(info.GetRole()),
	}
}

// ToUserInfoFromMessage to user info from message
func ToUserInfoFromMessage(info *model.UserInfoInMessage) *model.UserInfo {
	if info == nil {
		return nil
	}

	return &model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  info.Role,
	}
}

// ToLogUserMessageFromService to log user message
func ToLogUserMessageFromService(user *model.User) *model.LogUserMessage {
	if user == nil {
		return nil
	}

	return &model.LogUserMessage{
		ID: user.ID,
		Info: model.LogUserInfoInMessage{
			Name:  user.Info.Name,
			Email: user.Info.Email,
			Role:  user.Info.Role,
		},
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: func(updatedAt sql.NullTime) *time.Time {
			if updatedAt.Valid {
				return &updatedAt.Time
			}
			return nil
		}(user.UpdatedAt),
	}
}

// ToLogUpdateUserMessageFromService convert to log update user message
func ToLogUpdateUserMessageFromService(userID int64, info *model.UpdateUserInfo) *model.LogUpdateUserMessage {
	var infoInMessage *model.LogUpdateUserInfoInMessage
	if info != nil {
		infoInMessage = &model.LogUpdateUserInfoInMessage{
			Name: info.Name,
			Role: info.Role,
		}
	}
	return &model.LogUpdateUserMessage{
		ID:   userID,
		Info: infoInMessage,
	}
}
