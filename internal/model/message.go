package model

import (
	"time"
)

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

// LogUserMessage represents the log user in message
type LogUserMessage struct {
	ID        int64                `json:"id"`
	Info      LogUserInfoInMessage `json:"info"`
	Password  string               `json:"password"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt *time.Time           `json:"updatedAt"`
}

// LogUserInfoInMessage represents user info model
type LogUserInfoInMessage struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  int32  `json:"role"`
}

// LogUpdateUserMessage log update user message
type LogUpdateUserMessage struct {
	ID   int64                       `json:"id"`
	Info *LogUpdateUserInfoInMessage `json:"info"`
}

// LogUpdateUserInfoInMessage log update user info in message
type LogUpdateUserInfoInMessage struct {
	Name *string `json:"name"`
	Role *int32  `json:"role"`
}
