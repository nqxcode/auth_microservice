package model

import (
	"time"
)

// AccessibleRole represents the model of accessible role
type AccessibleRole struct {
	ID              int64
	Role            string
	EndpointAddress string
	CreatedAt       time.Time
}
