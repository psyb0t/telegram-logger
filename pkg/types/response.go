package types

// Response is the struct representing the body of the HTTP response
type Response struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}
