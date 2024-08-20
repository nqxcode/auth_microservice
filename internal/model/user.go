package model

import (
	"database/sql"
	"time"
)

// User represents the user model
type User struct {
	ID        int64
	Info      UserInfo
	Password  string `json:"-"`
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// UserInfo user info model
type UserInfo struct {
	Name  string
	Email string
	Role  int32
}

// UpdateUserInfo model ot update user info
type UpdateUserInfo struct {
	Name *string
	Role *int32
}
