package types

// Request is the struct representing the HTTP request body
type Request struct {
	Caller    string                 `json:"caller"`
	Time      string                 `json:"time"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Error     string                 `json:"error"`
	RequestID string                 `json:"requestID"`
	TraceID   string                 `json:"traceID"`
	SpanID    string                 `json:"spanID"`
	Data      map[string]interface{} `json:"data"`
}
