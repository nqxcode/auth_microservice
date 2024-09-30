package model

import (
	"time"
)

// AccessibleRole repository model
type AccessibleRole struct {
	ID              int64 `db:"accessible_role_id"`
	Role            string
	EndpointAddress string
	CreatedAt       time.Time
}
