package model

type Response struct {
	Message string      `json:"message"`
	TraceID string      `json:"trace_id"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
