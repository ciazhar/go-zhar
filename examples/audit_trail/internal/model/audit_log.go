package model

type AuditLog struct {
	EventType string `json:"event_type"`
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
	Payload   string `json:"payload"`
}
