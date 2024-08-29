package model

import (
	"database/sql"
	"time"
)

// User represents the user model
type User struct {
	ID        int64
	Info      UserInfo
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// GetID get user id
func (u User) GetID() int64 {
	return u.ID
}

// UserInfo represents user info model
type UserInfo struct {
	Name  string
	Email string
	Role  int32
}

// UpdateUserInfo represents model ot update user info
type UpdateUserInfo struct {
	Name *string
	Role *int32
}

// UserMessage represents user info in message
type UserMessage struct {
	Info            UserInfoInMessage `json:"info"`
	Password        string            `json:"password"`
	PasswordConfirm string            `json:"passwordConfirm"`
}

// UserInfoInMessage represents user info in message
type UserInfoInMessage struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  int32  `json:"role"`
}
