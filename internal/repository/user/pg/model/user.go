package model

import (
	"database/sql"
	"time"
)

// User user repository model
type User struct {
	ID        int64 `db:"user_id"`
	Name      string
	Email     string
	Password  string
	Role      int32
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
