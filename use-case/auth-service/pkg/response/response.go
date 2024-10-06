package response

type Response struct {
	RequestID string `json:"requestID,omitempty"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
	Data      any    `json:"data,omitempty"`
}
