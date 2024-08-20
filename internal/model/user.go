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

type UserInfo struct {
	Name  string
	Email string
	Role  int32
}

type UpdateUserInfo struct {
	Name *string
	Role *int32
}
