package dto

// Response represents a generic API response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(message string) *Response {
	return &Response{
		Success: true,
		Message: message,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(message string) *Response {
	return &Response{
		Success: false,
		Message: message,
	}
}
