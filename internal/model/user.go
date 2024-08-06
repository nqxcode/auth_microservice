package model

import (
	"database/sql"
	"time"
)

// User represents the user model
type User struct {
	ID        int64
	Name      string
	Email     string
	Role      int32
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
