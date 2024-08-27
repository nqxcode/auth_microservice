package model

import (
	"time"
)

// Log log repository model
type Log struct {
	ID        int64 `db:"log_id"`
	IP        string
	Message   string
	Payload   any
	CreatedAt time.Time
}
