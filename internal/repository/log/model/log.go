package model

import (
	"time"
)

type Log struct {
	ID        int64 `db:"log_id"`
	Ip        string
	Message   string
	Payload   any
	CreatedAt time.Time
}
