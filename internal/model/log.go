package model

// Log represents the model of log
type Log struct {
	Message string
	Payload any
	IP      string
}

// LogMessage represents log in message
type LogMessage struct {
	Message string `json:"message"`
	Payload any    `json:"payload"`
	IP      string `json:"ip"`
}
